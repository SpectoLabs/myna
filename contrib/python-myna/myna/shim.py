import tempfile
import os
import shutil

TEMPLATE = """
#!/bin/bash -e

unset CAPTURE
myna {application} "\$@"
"""

def _add_shim_directory_to_path(shim_dir):
    os.environ['PATH'] = shim_dir + os.path.pathsep + os.environ['PATH']

def _remove_shim_directory_from_path(shim_dir):
    if os.environ['PATH'].startswith(shim_dir + os.path.pathsep):
        os.environ['PATH'] = os.environ['PATH'][len(shim_dir + os.path.pathsep):]

def setup_shim_for(application):
    shim_dir = tempfile.mkdtemp(suffix='myna', prefix='tmp')
    shim_path = os.path.join(shim_dir, application)
    with open(shim_path, 'w') as f:
        f.write(TEMPLATE.format(application=application))
    os.chmod(shim_path, 0o750)
    _add_shim_directory_to_path(shim_dir)
    return shim_dir

def teardown_shim_dir(shim_dir):
    shutil.rmtree(shim_dir)
    _remove_shim_directory_from_path(shim_dir)

setup_shim_for('echo')
