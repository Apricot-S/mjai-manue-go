import subprocess

proc = subprocess.Popen(
    ["./mjai-manue", "--pipe"],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    text=True,
    bufsize=1,
)

try:
    while True:
        user_input = input()
        if user_input:
            proc.stdin.write(user_input)
            proc.stdin.flush()

        output = proc.stdout.readline()
        if output:
            print(output.strip())
except KeyboardInterrupt:
    proc.terminate()
    proc.wait()
