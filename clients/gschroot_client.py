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

	def _message(self, method, *params):
		msg = {"id": str(ObjectId()), "method": method, "params": list(params)}
		return json.dumps(msg).encode()

	def get_task_status(self):
		encoded_msg = self._message("GsServer.GetTaskStatus", {})
		self.socket.send(encoded_msg)
		chunk = self.socket.recv(2048)
		return chunk

	def get_task_stdout(self):
		encoded_msg = self._message("GsServer.GetTaskStdLogPath", {})
		self.socket.send(encoded_msg)
		chunk = self.socket.recv(2048)
		return chunk

	def get_task_stderr(self):
		pass