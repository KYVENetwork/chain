<!--
order: 4
-->

# BeginBlock

The `x/team` module begin-block hook handles distribution of inflation rewards for the team. It calculates
the current share of vested $KYVE in the team module and takes this share from the current block rewards and
transfers them to inflation reward wallet controlled by the team authority. Those $KYVE are intended as
reward for early team members since vesting starts earlier for early team members.
