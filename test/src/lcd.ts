import KyveSDK, { constants, KyveLCDClientType } from "@kyve/sdk";
import { JsonSchemaGenerator } from "typescript-json-schema/typescript-json-schema";
import { createValidator } from "./helpers/helper";
import { NETWORK } from "./helpers/constants";
import {
  ADDRESS_ALICE,
  ADDRESS_CHARLIE,
  alice,
  bob,
  charlie,
} from "./helpers/accounts";
import BigNumber from "bignumber.js";
const PATH_TO_QUERY_TYPES =
  "./node_modules/@kyve/proto/dist/proto-res/kyve/registry/v1beta1/query";

const TEST_HEIGHT = "0";

export const lcd = () => {
  let lcdClient: KyveLCDClientType;
  let typeQuerySchemas: JsonSchemaGenerator;
  let validate: Function;
  beforeAll(async () => {
    const sdk = new KyveSDK(NETWORK);
    lcdClient = await sdk.createLCDClient();
    const result = createValidator([PATH_TO_QUERY_TYPES]);
    validate = result.validate;
    typeQuerySchemas = result.typeQuerySchemas;

    const amount = new BigNumber(80)
      .multipliedBy(10 ** constants.KYVE_DECIMALS)
      .toString();
    //preparing data before lcd tests
    //found a pool
    await alice.client.kyve.v1beta1.base
      .fundPool({
        amount,
        id: "0",
      })
      .then((tx) => tx.execute());
    //stake pool
    await charlie.client.kyve.v1beta1.base
      .stakePool({
        id: "0",
        amount,
      })
      .then((tx) => tx.execute());
    //fund pool
    await alice.client.kyve.v1beta1.base
      .fundPool({
        id: "0",
        amount,
      })
      .then((tx) => tx.execute());

    //delegate
    await bob.client.kyve.v1beta1.base
      .delegatePool({
        id: "0",
        amount,
        staker: alice.client.account.address,
      })
      .then((tx) => tx.execute());
    //undelegate
    await bob.client.kyve.v1beta1.base.undelegatePool({
      id: "0",
      staker: ADDRESS_ALICE,
      amount,
    });
  });

  test("Query <params>", async () => {
    const result = await lcdClient.kyve.registry.v1beta1.params();
    const schema = typeQuerySchemas.getSchemaForSymbol("QueryParamsResponse");
    const validationResult = validate(schema, result);
    expect(validationResult.valid).toBeTruthy();
  });

  test("Query <pools> and <pool> by id", async () => {
    const poolsResponse = await lcdClient.kyve.registry.v1beta1.pools();
    const schema = typeQuerySchemas.getSchemaForSymbol("QueryPoolsResponse");
    //do not test pagination property
    delete schema.properties?.pagination;
    delete poolsResponse.pagination;
    const vResult = validate(schema, poolsResponse);
    expect(vResult.valid).toBeTruthy();
    // jest doesn't support nested generative test, needs a solution how to split into separate test cases
    // maybe another test runner?
    for (let pool of poolsResponse.pools) {
      const poolsResponse = await lcdClient.kyve.registry.v1beta1.pool({
        id: pool.id,
      });
      const schema = typeQuerySchemas.getSchemaForSymbol("QueryPoolResponse");
      const vResult = validate(schema, poolsResponse);
      expect(vResult.valid).toBeTruthy();
    }
  });

  test("Query <fundersList>", async () => {
    const poolsResponse = await lcdClient.kyve.registry.v1beta1.pools();
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryFundersListResponse"
    );
    for (let pool of poolsResponse.pools) {
      const poolsResponse = await lcdClient.kyve.registry.v1beta1.fundersList({
        pool_id: pool.id,
      });
      const vResult = validate(schema, poolsResponse);
      expect(vResult.valid).toBeTruthy();
    }
  });
  //
  test("Query <funder>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const founders = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: pool.pools[0].id,
    });
    const founder = await lcdClient.kyve.registry.v1beta1.funder({
      pool_id: pool.pools[0].id,
      funder: founders.funders[0].account,
    });
    const schema = typeQuerySchemas.getSchemaForSymbol("QueryFunderResponse");
    const vResult = validate(schema, founder);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <stakersList>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakers = await lcdClient.kyve.registry.v1beta1.stakersList({
      pool_id: pool.pools[0].id,
      status: 1,
    });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryStakersListResponse"
    );
    //do not test pagination property
    delete schema.properties?.pagination;
    delete stakers.pagination;
    const vResult = validate(schema, stakers);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <staker>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const stakerResponse = await lcdClient.kyve.registry.v1beta1.staker({
      pool_id: pool.pools[0].id,
      staker: stakersListResponse.stakers[0].staker,
    });
    const schema = typeQuerySchemas.getSchemaForSymbol("QueryStakerResponse");
    const vResult = validate(schema, stakerResponse);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <canPropose>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const canProposeRes = await lcdClient.kyve.registry.v1beta1.canPropose({
      pool_id: pool.pools[0].id,
      proposer: stakersListResponse.stakers[0].staker,
      from_height: TEST_HEIGHT,
    });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryCanProposeResponse"
    );
    const vResult = validate(schema, canProposeRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <stakeInfo>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const stakeInfoRes = await lcdClient.kyve.registry.v1beta1.stakeInfo({
      pool_id: pool.pools[0].id,
      staker: stakersListResponse.stakers[0].staker,
    });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryStakeInfoResponse"
    );
    const vResult = validate(schema, stakeInfoRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountAssets>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const accountAssetsRes =
      await lcdClient.kyve.registry.v1beta1.accountAssets({
        address: stakersListResponse.stakers[0].account,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountAssetsResponse"
    );
    const vResult = validate(schema, accountAssetsRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountFundedList>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const foundersRes = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: pool.pools[0].id,
    });
    const accountFundedListRes =
      await lcdClient.kyve.registry.v1beta1.accountFundedList({
        address: foundersRes.funders[0].account,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountFundedListResponse"
    );
    // do not test pagination property
    delete schema.properties?.pagination;
    delete accountFundedListRes.pagination;

    const vResult = validate(schema, accountFundedListRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountStakedList>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const accountStakedListRes =
      await lcdClient.kyve.registry.v1beta1.accountStakedList({
        address: stakersListResponse.stakers[0].account,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountStakedListResponse"
    );
    // do not test pagination property
    delete schema.properties?.pagination;
    delete accountStakedListRes.pagination;
    const vResult = validate(schema, accountStakedListRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountDelegationList>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const accountDelegationListRes =
      await lcdClient.kyve.registry.v1beta1.accountDelegationList({
        address: stakersListResponse.stakers[0].account,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountDelegationListResponse"
    );
    // do not test pagination property
    delete schema.properties?.pagination;
    delete accountDelegationListRes.pagination;
    const vResult = validate(schema, accountDelegationListRes);
    expect(vResult.valid).toBeTruthy();
  });
  test("Query <delegatorsByPoolAndStaker>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const delegatorsRes =
      await lcdClient.kyve.registry.v1beta1.delegatorsByPoolAndStaker({
        pool_id: pool.pools[0].id,
        staker: stakersListResponse.stakers[0].staker,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryDelegatorsByPoolAndStakerResponse"
    );
    //do not test pagination property
    delete schema.properties?.pagination;
    delete delegatorsRes.pagination;
    const vResult = validate(schema, delegatorsRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <delegator>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const delegatorsRes =
      await lcdClient.kyve.registry.v1beta1.delegatorsByPoolAndStaker({
        pool_id: pool.pools[0].id,
        staker: ADDRESS_ALICE,
      });

    const delegatorResponse = await lcdClient.kyve.registry.v1beta1.delegator({
      pool_id: pool.pools[0].id,
      staker: stakersListResponse.stakers[0].staker,
      delegator: delegatorsRes.delegators[0].delegator,
    });

    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryDelegatorResponse"
    );

    const vResult = validate(schema, delegatorResponse);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <stakersByPoolAndDelegator>", async () => {
    const pool = await lcdClient.kyve.registry.v1beta1.pools({
      pagination: { limit: "1" },
    });
    const stakersListResponse =
      await lcdClient.kyve.registry.v1beta1.stakersList({
        pool_id: pool.pools[0].id,
        status: 1,
      });
    const delegatorsRes =
      await lcdClient.kyve.registry.v1beta1.delegatorsByPoolAndStaker({
        pool_id: pool.pools[0].id,
        staker: ADDRESS_ALICE,
      });
    const stakersByPoolAndDelegatorRes =
      await lcdClient.kyve.registry.v1beta1.stakersByPoolAndDelegator({
        pool_id: pool.pools[0].id,
        delegator: delegatorsRes.delegators[0].delegator,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryStakersByPoolAndDelegatorResponse"
    );
    //do not test pagination property
    delete schema.properties?.pagination;
    delete stakersByPoolAndDelegatorRes.pagination;
    const vResult = validate(schema, stakersByPoolAndDelegatorRes);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountStakingUnbondingsRequest>", async () => {
    const accountStakingUnbondingResponse =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_ALICE,
      });
    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountStakingUnbondingsResponse"
    );
    delete schema.properties?.pagination;
    delete accountStakingUnbondingResponse.pagination;
    const vResult = validate(schema, accountStakingUnbondingResponse);
    expect(vResult.valid).toBeTruthy();
  });

  test("Query <accountDelegationUnbondings>", async () => {
    const accountDelegationUnbondingRes =
      await lcdClient.kyve.registry.v1beta1.accountDelegationUnbondings({
        address: bob.client.account.address,
      });

    const schema = typeQuerySchemas.getSchemaForSymbol(
      "QueryAccountDelegationUnbondingsResponse"
    );
    delete schema.properties?.pagination;
    delete accountDelegationUnbondingRes?.pagination;
    const vResult = validate(schema, accountDelegationUnbondingRes);
    expect(vResult.valid).toBeTruthy();
  });
  //todo: add query tests
  // - canVote
  // - proposals
  // - proposal
  // - proposal by hight
};
