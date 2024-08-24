import subprocess
import time

subprocess.Popen(['curl', '-v', 'stefandejager.nl'])

while True:
    subprocess.Popen(['echo', 'Hello from subprocess'])
    time.sleep(3000)


