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
		encoded_msg = json.dumps(msg).encode()
		self.socket.send(encoded_msg)
		chunk = self.socket.recv(2048)
		return chunk

	def get_task_status(self):
		chunk = self._message("GsServer.GetTaskStatus", {})
		return json.loads(chunk)

	def get_task_stdout(self):
		"""returns stdout from command"""
		#TODO: make it so the whole file doesn't need to get read into memory at once
		chunk = self._message("GsServer.GetTaskStdLogPath", {})
		log_path_dict = json.loads(chunk)
		fname = log_path_dict["result"]["StdoutPath"]
		with open(fname) as f:
			contents = f.read()
		return contents

	def get_task_stderr(self):
		"""returns stderr from command"""
		pass