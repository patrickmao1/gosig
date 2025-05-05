import geni.portal as portal
import geni.rspec.pg as rspec

NUM_NODES = 5

pc = portal.Context()
request = rspec.Request()

lan = request.LAN('lan0')
lan.best_effort = False

for i in range(NUM_NODES):
    name = "node{}".format(i)
    node = request.RawPC(name)
    iface = node.addInterface("if{}".format(i))
    iface.addAddress(rspec.IPv4Address("10.0.0.%d" % (10 + i), "255.255.255.0"))
    lan.addInterface(iface)

    cmd = """
sudo apt-get update
cd /users/patrickm
wget -P /tmp https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/go1.24.2.linux-amd64.tar.gz
echo "export node_index={0}" >> /users/patrickm/.bashrc
echo "export PATH=$PATH:/usr/local/go/bin" >> /users/patrickm/.bashrc
source /users/patrickm/.bashrc
git clone https://github.com/patrickmao1/gosig && cd gosig || exit
make install
    """.format(i)

    node.addService(rspec.Execute(shell="bash", command=cmd))

pc.printRequestRSpec(request)
