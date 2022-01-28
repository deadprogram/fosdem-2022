connect:
	cd ./demo/connect; tinygo flash -size short -target=arduino-nano33 .; cd ../..

webclient:
	cd ./demo/webclient; tinygo flash -size short -target=wioterminal .; cd ../..

mqttclient:
	cd ./demo/mqttclient; tinygo flash -size short -target=pyportal .; cd ../..

flightbadge:
	cd ./demo/flightbadge; tinygo flash -size short -target=pybadge .; cd ../..
