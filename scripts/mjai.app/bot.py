import subprocess

proc = subprocess.Popen(
    ["./mjai-manue", "--pipe"],
    bufsize=1,
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    text=True,
    encoding="utf-8",
)

# subprocess.Popen's stdin/stdout/stderr might be None depending on
# the pipe settings
# See: https://docs.python.org/3.13/library/subprocess.html#subprocess.Popen.stdin
assert proc.stdin is not None
assert proc.stdout is not None
assert proc.stderr is not None

try:
    while True:
        input_ = input()
        if input_.strip() == "":
            # Workaround
            continue
        proc.stdin.write(input_)
        proc.stdin.flush()

        output = proc.stdout.readline()
        if output:
            print(output.strip())
except KeyboardInterrupt:
    proc.terminate()
    proc.wait()
