# yolo

An experimental CLI for tweaking existing Cog models and deploying them to Replicate really fast.

- No Docker required
- No Cog required
- No GPU required

## DISCLAIMER

⚠️ This is a tool for power users that has many rough edges. ⚠️

## Usage

### Install on mac/linux

Fetch the precompiled Go binary from GitHub Releases:

    sudo curl -o /usr/local/bin/yolo -L "https://github.com/replicate/yolo/releases/latest/download/yolo_$(uname -s)_$(uname -m)"
    sudo chmod +x /usr/local/bin/yolo

Alternatively can build from source (Golang required):

    go build && sudo cp yolo /usr/local/bin

### Get your Replicate token

You can use **either** your REPLICATE_API_TOKEN or your COG_TOKEN.

By using your REPLICATE_API_TOKEN, we can access the API and PUSH models to your account.

#### Replicate API Token

Visit https://replicate.com/account/api-tokens and copy your token.

    export REPLICATE_API_TOKEN=r8_...

#### Cog Token

Visit https://replicate.com/auth/token and copy your token.

    export COG_TOKEN=4b212....

### Modify a model (e.g. SDXL)

Grab the code by cloning the repo

   git clone https://github.com/replicate/cog-sdxl.git

### Find an existing version to modify

Visit https://replicate.com/stability-ai/sdxl/api and find the docker image name:

    r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41

This is going to be your "base" for your tweaked model.  You can think 
of the process as adding your changes on top of this model, as that is
what happens under the hood.  A new layer is added with whatever files
you specify.

### Create a model

There is no Replicate API for model creation, so you must create a model on the website:

    https://replicate.com/create

### Make and push your changes

If you are NOT changing the schema (inputs/outputs), you run this:

    yolo push \
    --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 \
    --dest r8.im/anotherjesse/my-awesome-changes \
    list_of_files_to_send

If you are changing the schema

    yolo push \
    --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 \
    --dest r8.im/anotherjesse/my-awesome-changes \
    --ast path_to_predictor.py \
    path_to_predictor.py other_files_here.py



