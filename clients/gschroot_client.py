from bson import ObjectId
import json
import subprocess
import socket
import multiprocessing
import time
import os
import signal
SOCK_DIR = "/tmp/gschroot_socks"
class Gschroot(object):
	def __init__(self, url, cmd, name, retries = 3):
		self.proc = subprocess.Popen("main -url=%s -cmd='%s' -name=%s" % (url, cmd, name), shell=True, stdout = subprocess.PIPE, preexec_fn=os.setsid)
		self.name = name
		self.sock_path = SOCK_DIR + "/" + name
		self.socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
		time.sleep(.1) #work around to make sure (more like give better chance) process has started
		for x in xrange(retries):
			try:
				self.socket.connect(self.sock_path)
				break
			except Exception as inst:
				if x == retries - 1:
					raise 

	def __del__(self):
		os.killpg(self.proc.pid, signal.SIGTERM)

	def get_task_status(self):
		self.socket.send(json.dumps({"id": str(ObjectId()), "method": "GsServer.GetTaskStatus", "params": [{"pid": self.proc.pid}]}).encode())
		chunk = self.socket.recv(2048)
		print chunk
		#output = subprocess.check_output("main -query='%s'" % (self.name), shell=True)
		#print output
		#return json.loads(output)

	def get_task_stdout(self):
		pass

	def get_task_stderr(self):
		pass