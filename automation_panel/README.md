# Automation Panel

Build container image using `bash build.sh <version>` command. The version is a semver tag like `v0.5.0`. The tag is used to pin point the container image in the clab file.

Note: If you change something in the panel, you need to build a new panel container with `bash build.sh <version>` and upload to GH.
For testing you can push to your own account and change the path in the topology to it, or you can create the panel in your local docker. Refer to the instruction inside the script file.

