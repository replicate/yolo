# YOLO

Goals - from my mac - without docker I want to tweak SDXL and deploy in replicate in seconds

1. update predict.py doing doing control reference
 - update t2i to return controls
2. update params
 - parse ast to 
3. allow cog weights to be set?
 - not sure if this is something completely different - since combining code changes and hot swaps is confusing
3.b. help manage weights
 - api to upload weights to replicate-weights / trigger sync to all the clouds
4. basic checks
 - check cog.yaml hasn't changed
 - check git shas?
5. github action
 - do this as an action?