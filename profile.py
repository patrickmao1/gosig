import geni.portal as portal
import geni.rspec.pg as rspec

# 1. Create a Portal (for params) and RSpec Request
pc = portal.Context()
request = rspec.Request()

# 2. Add nodes
node = request.RawPC('node1')
# (Optional) set node.disk_image, node.hardware_type, etc.
iface = node.addInterface('if1')
# 3. Add networks (LAN)
lan = request.LAN('lan0')
lan.addInterface(iface)

# 4. Add services or commands to run on each node
node.addService(rspec.Execute(shell='bash', command='your_startup_script.sh'))

# 5. Print the generated RSpec
portal.Context().printRequestRSpec(request)
