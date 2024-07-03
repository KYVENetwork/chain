# Contributing


Thank you for considering to contribute to this project. KYVE is an L1 based
on Cosmos SDK and CometBFT. We mostly follow their principles and design
architectures.

## Overview

- The latest state of development is on `main`.
- `main` must always pass `make all ENV=mainnet`.
- Releases can be found in `/release/*`.
- Everything must be covered by tests. We have a very extensive test-suite
  and use triple-A testing (Arrange, Act, Assert).

## Creating a Pull Request

- Check out the latest state from main and always keep the PR in sync with main.
- Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#specification).
- Only one feature per pull request.
- Write an entry for the Changelog.
- Write tests covering 100% of your modified code.
- The command `make all ENV=mainnet` must pass. 

## Coding Guidelines

- Write readable and maintainable code. `Premature Optimization Is the Root of All Evil`.
  Concentrate on clean interfaces first and only optimize for performance if it is needed.
- The keeper directory is structured the following:
  - `getters_*`-files only interact with the KV-Store. All methods always succeed
    and do not return errors. This is the only place where methods are allowed to 
    write to the KV-Store. Also, all aggregation variables are updated here.
  - `logic_*`-files handle complex tasks and are encouraged to emit events and
    call the getters functions. 
  - `msg_server_*`-files are the entry point for message handling. This file
    should be very clean to read and outsource most of the part to the logic files.
    One should immediately understand the flow by just reading the function names
    which are called while handling the message.

## Legal

You agree that your contribution is licenced under the MIT Licence and all
ownership is handed over the authors named in `LICENSE`.
