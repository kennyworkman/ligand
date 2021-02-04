ligand 🔌
----

A local entrypoint to cloud compute, built for bioinformaticians.

_Note that this repo has not had an official release yet and is still experimental!_

<img src="https://github.com/latchai/latch-docs/raw/main/static/images/latch-cropped-demo.gif" width="700" />

### Quickstart

```
pip install ligand
```
Then in any script:

```python
import ligand

# Make sure the init is directly beneath your imports!

ligand.init()

# And above the code you want to run.
```

Thats it. Told you it was easy...

Now just run your script like you would normally, ie. `python train.py` and it
will run on the cloud.


### Features

Busy building, coming soon...

### Acknowledgements

Architectural decisions and subroutines, especially with respect to interaction between a
compiled daemon and a python SDK, were inspired or borrowed from the
[replicate.ai](https://replicate.ai) project.
