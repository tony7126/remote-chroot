from gschroot_client import Gschroot
import unittest

class TestGschrootClient(unittest.TestCase):
	def test_gschroot(self):
		gs = Gschroot("http://localhost:8000/rootfs.tar", cmd="ls", name = "testing")
		print gs.get_task_status()
		print gs.get_task_stdout()


if __name__ == "__main__":
	unittest.main()