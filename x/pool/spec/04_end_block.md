<!--
order: 4
-->

# EndBlock

EndBlock is used to determine if a scheduled runtime upgrade needs to be performed based on the
provided upgrade time. If an upgrade is scheduled and the scheduled time is reached _end_block_ will copy over
the upgrade details to the actual pool version and pauses the pool for the specified duration. After the end of the
duration is reached _end_block_ again unpauses the pool, finishing the runtime upgrade.