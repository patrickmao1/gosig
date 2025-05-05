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
echo $(pwd)
sudo apt-get update
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo 'export node_index={0}' >> ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
git clone https://github.com/patrickmao1/gosig && cd gosig
make install
""".format(i, name)

    node.addService(rspec.Execute(shell="bash", command=cmd))

pc.printRequestRSpec(request)