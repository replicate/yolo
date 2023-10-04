# YOLO

Goals - from my mac - without docker I want to tweak existing models and deploy in replicate in seconds.

This will not rebuild changes in cog.yaml (no new dependencies), ...

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

# Usage

## Build on mac/linux

    go build

This builds `./yolo`, you can install by copying into your path.

    sudo cp yolo /usr/local/bin

# Get your replicate CLI auth token

Visit https://replicate.com/auth/token and copy your token.

    export COG_TOKEN=4b212....

## modify an model (SDXL)

Grab the code by cloning the repo

   git clone https://github.com/replicate/cog-sdxl.git

## find an existing version to modify

Visit https://replicate.com/stability-ai/sdxl/api find the docker image name

    r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41

This is going to be your "base" for your tweaked model.  You can think 
of the process as adding your changes on top of this models, as that is
what happens under the hood.  A new layer is added with whatever files
you specify.

## Create a model

Since there is no API for model creation, you must already have created a model

    https://replicate.com/create

## make and push your changes

If you are NOT changing the schema (inputs/outputs), you run this:

    yolo push --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 --dest r8.im/anotherjesse/my-awesome-changes list_of_files_to_send

If you are changing the schema

    yolo push --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 --dest r8.im/anotherjesse/my-awesome-changes --ast path_to_predictor.py path_to_predictor.py other_files_here.py



