import subprocess


print("Hello, World!")


# Python 2
#subprocess.Popen("ssh {user}@{host} {cmd}".format(user=user, host=host, cmd='ls -l'), shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE).communicate()

# Python 3
#subprocess.Popen(f"ssh {user}@{host} {cmd}", shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE).communicate()


#sshpass -p {password} ssh {user}@{ip}
#f"echo {password} | ssh {user}@{host} {cmd}"


#subprocess.Popen(f"ssh -p {password} {user}@{host} {cmd}", shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE).communicate()