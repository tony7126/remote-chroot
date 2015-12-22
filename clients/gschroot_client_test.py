from gschroot_client import Gschroot
import unittest

class TestGschrootClient(unittest.TestCase):
	def test_gschroot(self):
		gs = Gschroot("http://localhost:8000/rootfs.tar", cmd="sleep 3", name = "testing")
		gs.get_task_status()


if __name__ == "__main__":
	unittest.main()