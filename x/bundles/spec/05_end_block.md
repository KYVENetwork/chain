<!--
order: 5
-->

# EndBlock

EndBlock is used to determine if the uploader did not
submit his bundle proposal in a predefined timeout. The penalty
for not uploading in time is a point. If a participant reaches
a certain number of points the participant receives a timeout slash
and gets removed from the storage pool.

To prevent that the uploader should always upload a bundle proposal.
If he can not do that for whatever reason the uploader should skip
his uploader role, indicating he is not offline.