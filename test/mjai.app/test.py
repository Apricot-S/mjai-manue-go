import random
import sys

import mjai

if len(sys.argv) != 3:
    print("Usage: test.py <logs_dir> <archive_path>")
    sys.exit(1)

logs_dir = sys.argv[1]
archive_path = sys.argv[2]

submissions = [
    archive_path,
    archive_path,
    archive_path,
    archive_path,
]

mjai.Simulator(
    submissions,
    logs_dir=logs_dir,
    seed=(random.randint(0, sys.maxsize), random.randint(0, sys.maxsize)),
    timeout=10,
).run()
