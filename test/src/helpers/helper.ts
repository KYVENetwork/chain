// optionally pass argument to schema generator
import * as TJS from "typescript-json-schema";
import { Definition } from "typescript-json-schema";
import { ErrorObject } from "ajv/lib/types";
import { JsonSchemaGenerator } from "typescript-json-schema/typescript-json-schema";
import AJV from "ajv";

const settings: TJS.PartialArgs = {
  // strictNullChecks: false,
  required: true,
  noExtraProps: true,
};

// optionally pass ts compiler options
const compilerOptions: TJS.CompilerOptions = {
  strictNullChecks: false,
  additionalProperties: false,
};
const removeNull = (obj: any) => {
  Object.keys(obj).forEach(
    (key) =>
      (obj[key] && typeof obj[key] === "object" && removeNull(obj[key])) ||
      ((obj[key] === undefined || obj[key] === null) && delete obj[key])
  );
  return obj;
};
export function createValidator(pathToTypes: string[]) {
  const program = TJS.getProgramFromFiles([...pathToTypes], compilerOptions);
  const typeQuerySchemas = TJS.buildGenerator(
    program,
    settings
  ) as unknown as JsonSchemaGenerator;
  const ajv = new AJV({ strict: true });
  if (!typeQuerySchemas) {
    throw new Error("Can't find query type to generate JSON schema ");
  }

  function validate(
    schema: Definition,
    data: any
  ): { valid: boolean; errors: ErrorObject[] | null } {
    const validate = ajv.compile(schema);

    validate(removeNull(data));
    if (validate.errors)
      return {
        valid: false,
        errors: validate.errors,
      };
    return {
      valid: true,
      errors: null,
    };
  }
  return {
    typeQuerySchemas,
    validate,
  };
}
