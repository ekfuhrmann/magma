"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import xml.etree.ElementTree as ET
from datetime import datetime, timezone, timedelta
from unittest import TestCase, mock
from unittest.mock import Mock, patch
from spyne import MethodContext
from spyne.server import ServerBase
from magma.enodebd.devices.baicells import BaicellsHandler
from magma.enodebd.state_machines.enb_acs_pointer import StateMachinePointer
from magma.enodebd.stats_manager import StatsManager
from magma.enodebd.tr069 import models
from magma.enodebd.tr069.rpc_methods import AutoConfigServer
from magma.enodebd.tr069.spyne_mods import Tr069Application, Tr069Soap11


class Tr069Test(TestCase):
    """ Tests for the TR-069 server """
    acs_to_cpe_queue = None
    cpe_to_acs_queue = None

    def setUp(self):
        # Build the state_machine
        self.acs_state_machine = self._build_acs_state_machine()

        # Set up the ACS
        state_machine_pointer = StateMachinePointer(self.acs_state_machine)
        AutoConfigServer.set_state_machine_pointer(state_machine_pointer)

        def side_effect(*args, **_kwargs):
            msg = args[1]
            return msg

        self.p = patch.object(AutoConfigServer, '_handle_tr069_message',
                              Mock(side_effect=side_effect))
        self.p.start()

        self.app = Tr069Application([AutoConfigServer],
                                    models.CWMP_NS,
                                    in_protocol=Tr069Soap11(validator='soft'),
                                    out_protocol=Tr069Soap11())

    def _build_acs_state_machine(self) -> BaicellsHandler:
        # Build the state_machine
        stats_mgr = StatsManager()
        event_loop = asyncio.get_event_loop()
        mconfig = self._get_mconfig()
        service_config = self._get_service_config()
        with mock.patch('magma.common.service.MagmaService') as MockService:
            MockService.config = service_config
            MockService.mconfig = mconfig
            MockService.loop = event_loop
            acs_state_machine = BaicellsHandler(MockService, stats_mgr)
            return acs_state_machine

    def tearDown(self):
        self.p.stop()

    def _get_mconfig(self):
        return {
            "@type": "type.googleapis.com/magma.mconfig.EnodebD",
            "bandwidthMhz": 20,
            "specialSubframePattern": 7,
            "earfcndl": 44490,
            "logLevel": "INFO",
            "plmnidList": "00101",
            "pci": 260,
            "allowEnodebTransmit": False,
            "subframeAssignment": 2,
            "tac": 1
        },

    def _get_service_config(self):
        return {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
            "reboot_enodeb_on_mme_disconnected": True,
            "s1_interface": "eth1",
        }

    def test_parse_inform(self):
        """
        Test that example Inform RPC call can be parsed correctly
        """
        # Example TR-069 CPE->ACS RPC call. Copied from:
        # http://djuro82.blogspot.com/2011/05/tr-069-cpe-provisioning.html
        cpe_string = b'''
            <soapenv:Envelope soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
                <soapenv:Header>
                    <cwmp:ID soapenv:mustUnderstand="1">0_THOM_TR69_ID</cwmp:ID>
                </soapenv:Header>
                <soapenv:Body>
                    <cwmp:Inform>
                        <DeviceId>
                            <Manufacturer>THOMSON</Manufacturer>
                            <OUI>00147F</OUI>
                            <ProductClass>SpeedTouch 780</ProductClass>
                            <SerialNumber>CP0611JTLNW</SerialNumber>
                        </DeviceId>
                        <Event soap:arrayType="cwmp:EventStruct[04]">
                            <EventStruct>
                                <EventCode>0 BOOTSTRAP</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>1 BOOT</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>2 PERIODIC</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>4 VALUE CHANGE</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                        </Event>
                        <MaxEnvelopes>2</MaxEnvelopes>
                        <CurrentTime>1970-01-01T00:01:09Z</CurrentTime>
                        <RetryCount>05</RetryCount>
                        <ParameterList soap:arrayType="cwmp:ParameterValueStruct[12]">
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceSummary</Name>
                                <Value xsi:type="xsd:string">
                                    InternetGatewayDevice:1.1[] (Baseline:1, EthernetLAN:1, ADSLWAN:1, Bridging:1, Time:1, WiFiLAN:1)</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
                                <Value xsi:type="xsd:string">1.1</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
                                <Value xsi:type="xsd:string">BANT-R</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
                                <Value xsi:type="xsd:string">6.2.35.0</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
                                <Value xsi:type="xsd:string"></Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Name</Name>
                                <Value xsi:type="xsd:string">MyCompanyName</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Version</Name>
                                <Value xsi:type="xsd:string"></Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Date</Name>
                                <Value xsi:type="xsd:dateTime">0001-01-01T00:00:00</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceInfo.VendorConfigFile.1.Description</Name>
                                <Value xsi:type="xsd:string">MyCompanyName</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name>
                                <Value xsi:type="xsd:string">http://10.127.129.205:51005/</Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.ManagementServer.ParameterKey</Name>
                                <Value xsi:type="xsd:string"></Value>
                            </ParameterValueStruct>
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name>
                                <Value xsi:type="xsd:string">10.127.129.205</Value>
                            </ParameterValueStruct>
                        </ParameterList>
                    </cwmp:Inform>
                </soapenv:Body>
            </soapenv:Envelope>
            '''

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        if ctx.in_error is not None:
            print('In error: %s' % ctx.in_error)
        self.assertEqual(ctx.in_error, None)

        server.get_in_object(ctx)

        self.assertEqual(ctx.in_object.DeviceId.OUI, '00147F')
        self.assertEqual(
            ctx.in_object.Event.EventStruct[0].EventCode, '0 BOOTSTRAP')
        self.assertEqual(
            ctx.in_object.Event.EventStruct[2].EventCode, '2 PERIODIC')
        self.assertEqual(ctx.in_object.MaxEnvelopes, 2)
        self.assertEqual(
            ctx.in_object.ParameterList.ParameterValueStruct[1].Name,
            'InternetGatewayDevice.DeviceInfo.SpecVersion')
        self.assertEqual(
            str(ctx.in_object.ParameterList.ParameterValueStruct[1].Value), '1.1')

    def test_parse_inform_cavium(self):
        """
        Test that example Inform RPC call can be parsed correctly from OC-LTE
        """
        cpe_string = b'''<?xml version="1.0" encoding="UTF-8"?>
        <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
          <SOAP-ENV:Header>
            <cwmp:ID SOAP-ENV:mustUnderstand="1">CPE_1002</cwmp:ID>
          </SOAP-ENV:Header>
          <SOAP-ENV:Body SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
            <cwmp:Inform>
              <DeviceId>
                <Manufacturer>Cavium, Inc.</Manufacturer>
                <OUI>000FB7</OUI>
                <ProductClass>Cavium eNB</ProductClass>
                <SerialNumber>10.18.104.79</SerialNumber>
              </DeviceId>
              <Event xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="cwmp:EventStruct[1]">
                <EventStruct>
                  <EventCode>0 BOOTSTRAP</EventCode>
                  <CommandKey></CommandKey>
                </EventStruct>
              </Event>
              <MaxEnvelopes>1</MaxEnvelopes>
              <CurrentTime>1970-01-02T00:01:05.021239+00:00</CurrentTime>
              <RetryCount>2</RetryCount>
              <ParameterList xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="cwmp:ParameterValueStruct[15]">
                <ParameterValueStruct>
                  <Name>Device.DeviceInfo.HardwareVersion</Name>
                  <Value xsi:type="xsd:string">1.0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.DeviceInfo.SoftwareVersion</Name>
                  <Value xsi:type="xsd:string">1.0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.DeviceInfo.AdditionalHardwareVersion</Name>
                  <Value xsi:type="xsd:string">1.0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.DeviceInfo.AdditionalSoftwareVersion</Name>
                  <Value xsi:type="xsd:string">1.0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.DeviceInfo.ProvisioningCode</Name>
                  <Value xsi:type="xsd:string">Cavium</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.ManagementServer.ParameterKey</Name>
                  <Value xsi:type="xsd:string"></Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.ManagementServer.ConnectionRequestURL</Name>
                  <Value xsi:type="xsd:string">http://192.88.99.253:8084/bucrhzjd</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.ManagementServer.UDPConnectionRequestAddress</Name>
                  <Value xsi:type="xsd:string"></Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.ManagementServer.NATDetected</Name>
                  <Value xsi:type="xsd:boolean">0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.PacketsReceived</Name>
                  <Value xsi:type="xsd:unsignedInt">0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.PacketsResponded</Name>
                  <Value xsi:type="xsd:unsignedInt">0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.BytesReceived</Name>
                  <Value xsi:type="xsd:unsignedInt">0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.BytesResponded</Name>
                  <Value xsi:type="xsd:unsignedInt">0</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.TimeFirstPacketReceived</Name>
                  <Value xsi:type="xsd:dateTime">1969-12-31T16:00:00.000000+00:00</Value>
                </ParameterValueStruct>
                <ParameterValueStruct>
                  <Name>Device.IP.Diagnostics.UDPEchoConfig.TimeLastPacketReceived</Name>
                  <Value xsi:type="xsd:dateTime">1969-12-31T16:00:00.000000+00:00</Value>
                </ParameterValueStruct>
              </ParameterList>
            </cwmp:Inform>
          </SOAP-ENV:Body>
        </SOAP-ENV:Envelope>
        '''

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        if ctx.in_error is not None:
            print('In error: %s' % ctx.in_error)
        self.assertEqual(ctx.in_error, None)

        server.get_in_object(ctx)

        self.assertEqual(ctx.in_object.DeviceId.OUI, '000FB7')
        self.assertEqual(
            ctx.in_object.Event.EventStruct[0].EventCode, '0 BOOTSTRAP')
        self.assertEqual(ctx.in_object.MaxEnvelopes, 1)
        self.assertEqual(
            ctx.in_object.ParameterList.ParameterValueStruct[1].Name,
            'Device.DeviceInfo.SoftwareVersion')
        self.assertEqual(
            str(ctx.in_object.ParameterList.ParameterValueStruct[1].Value), '1.0')

    def test_handle_transfer_complete(self):
        """
        Test that example TransferComplete RPC call can be parsed correctly, and
        response is correctly generated.
        """
        # Example TransferComplete CPE->ACS RPC request/response.
        # Manually created.
        cpe_string = b'''
            <soapenv:Envelope soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
                <soapenv:Header>
                    <cwmp:ID soapenv:mustUnderstand="1">1234</cwmp:ID>
                </soapenv:Header>
                <soapenv:Body>
                    <cwmp:TransferComplete>
                        <CommandKey>Downloading stuff</CommandKey>
                        <FaultStruct>
                            <FaultCode>0</FaultCode>
                            <FaultString></FaultString>
                        </FaultStruct>
                        <StartTime>2016-11-30T10:16:29Z</StartTime>
                        <CompleteTime>2016-11-30T10:17:05Z</CompleteTime>
                    </cwmp:TransferComplete>
                </soapenv:Body>
            </soapenv:Envelope>
            '''
        expected_acs_string = b'''
            <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
                <soapenv:Header>
                    <cwmp:ID soapenv:mustUnderstand="1">1234</cwmp:ID>
                </soapenv:Header>
                <soapenv:Body>
                    <cwmp:TransferCompleteResponse>
                    </cwmp:TransferCompleteResponse>
                </soapenv:Body>
            </soapenv:Envelope>
            '''

        self.p.stop()
        self.p.start()

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        if ctx.in_error is not None:
            print('In error: %s' % ctx.in_error)
        self.assertEqual(ctx.in_error, None)

        server.get_in_object(ctx)
        self.assertEqual(ctx.in_error, None)

        server.get_out_object(ctx)
        self.assertEqual(ctx.out_error, None)

        output_msg = ctx.out_object[0]
        self.assertEqual(type(output_msg), models.TransferComplete)
        self.assertEqual(output_msg.CommandKey, 'Downloading stuff')
        self.assertEqual(output_msg.FaultStruct.FaultCode, 0)
        self.assertEqual(output_msg.FaultStruct.FaultString, '')
        self.assertEqual(output_msg.StartTime,
                         datetime(2016, 11, 30, 10, 16, 29,
                                  tzinfo=timezone(timedelta(0))))
        self.assertEqual(output_msg.CompleteTime,
                         datetime(2016, 11, 30, 10, 17, 5,
                                  tzinfo=timezone(timedelta(0))))

        server.get_out_string(ctx)
        self.assertEqual(ctx.out_error, None)

        xml_tree = XmlTree()
        match = xml_tree.xml_compare(
            xml_tree.convert_string_to_tree(b''.join(ctx.out_string)),
            xml_tree.convert_string_to_tree(expected_acs_string))
        self.assertTrue(match)

    def test_parse_empty_http(self):
        """
        Test that empty HTTP message gets correctly mapped to 'EmptyHttp'
        function call
        """
        cpe_string = b''

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        if ctx.in_error is not None:
            print('In error: %s' % ctx.in_error)

        self.assertEqual(ctx.in_error, None)
        self.assertEqual(ctx.function, AutoConfigServer.empty_http)

    def test_generate_empty_http(self):
        """
        Test that empty HTTP message is generated when setting output message
        name to 'EmptyHttp'
        """
        cpe_string = b''

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        server.get_in_object(ctx)
        if ctx.in_error is not None:
            raise ctx.in_error

        server.get_out_object(ctx)
        if ctx.out_error is not None:
            raise ctx.out_error

        ctx.descriptor.out_message.Attributes.sub_name = 'EmptyHttp'
        ctx.out_object = [models.AcsToCpeRequests()]

        server.get_out_string(ctx)

        self.assertEqual(b''.join(ctx.out_string), b'')

    def test_generate_get_parameter_values_string(self):
        """
        Test that correct string is generated for SetParameterValues ACS->CPE
        request
        """
        # Example ACS->CPE RPC call. Copied from:
        # http://djuro82.blogspot.com/2011/05/tr-069-cpe-provisioning.html
        # Following edits made:
        # - Change header ID value from 'null0' to 'null', to match magma
        #   default ID
        expected_acs_string = b'''
        <soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
            <soapenv:Header>
                <cwmp:ID soapenv:mustUnderstand="1">null</cwmp:ID>
            </soapenv:Header>
            <soapenv:Body>
                <cwmp:GetParameterValues>
                    <ParameterNames soap:arrayType="xsd:string[1]">
                        <string>foo</string>
                    </ParameterNames>
                </cwmp:GetParameterValues>
            </soapenv:Body>
        </soapenv:Envelope>
        '''

        names = ['foo']
        request = models.GetParameterValues()
        request.ParameterNames = models.ParameterNames()
        request.ParameterNames.arrayType = 'xsd:string[%d]' \
            % len(names)
        request.ParameterNames.string = []
        for name in names:
            request.ParameterNames.string.append(name)

        request.ParameterKey = 'null'

        def side_effect(*_args, **_kwargs):
            return request

        AutoConfigServer.state_machine().handle_tr069_message = Mock(
            side_effect=side_effect)

        self.p.stop()
        self.p = patch.object(self.acs_state_machine, 'handle_tr069_message',
                              Mock(side_effect=side_effect))
        self.p.start()

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [b'']
        ctx, = server.generate_contexts(ctx)

        server.get_in_object(ctx)
        if ctx.in_error is not None:
            raise ctx.in_error

        server.get_out_object(ctx)
        if ctx.out_error is not None:
            raise ctx.out_error

        server.get_out_string(ctx)

        xml_tree = XmlTree()
        match = xml_tree.xml_compare(
            xml_tree.convert_string_to_tree(b''.join(ctx.out_string)),
            xml_tree.convert_string_to_tree(expected_acs_string))
        self.assertTrue(match)

    def test_generate_set_parameter_values_string(self):
        """
        Test that correct string is generated for SetParameterValues ACS->CPE
        request
        """
        # Example ACS->CPE RPC call. Copied from:
        # http://djuro82.blogspot.com/2011/05/tr-069-cpe-provisioning.html
        # Following edits made:
        # - Change header ID value from 'null0' to 'null', to match magma
        #   default ID
        expected_acs_string = b'''
        <soapenv:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
            <soapenv:Header>
                <cwmp:ID soapenv:mustUnderstand="1">null</cwmp:ID>
            </soapenv:Header>
            <soapenv:Body>
                <cwmp:SetParameterValues>
                    <ParameterList soap:arrayType="cwmp:ParameterValueStruct[4]">
                        <ParameterValueStruct>
                            <Name>InternetGatewayDevice.ManagementServer.PeriodicInformEnable</Name>
                            <Value xsi:type="xsd:boolean">1</Value>
                        </ParameterValueStruct>
                        <ParameterValueStruct>
                            <Name>InternetGatewayDevice.ManagementServer.ConnectionRequestUsername</Name>
                            <Value xsi:type="xsd:string">00147F-SpeedTouch780-CP0611JTLNW</Value>
                        </ParameterValueStruct>
                        <ParameterValueStruct>
                            <Name>InternetGatewayDevice.ManagementServer.ConnectionRequestPassword</Name>
                            <Value xsi:type="xsd:string">98ff55fb377bf724c625f60dec448646</Value>
                        </ParameterValueStruct>
                        <ParameterValueStruct>
                            <Name>InternetGatewayDevice.ManagementServer.PeriodicInformInterval</Name>
                            <Value xsi:type="xsd:unsignedInt">60</Value>
                        </ParameterValueStruct>
                    </ParameterList>
                    <ParameterKey>null</ParameterKey>
                </cwmp:SetParameterValues>
            </soapenv:Body>
        </soapenv:Envelope>
        '''

        request = models.SetParameterValues()

        request.ParameterList = \
            models.ParameterValueList(arrayType='cwmp:ParameterValueStruct[4]')
        request.ParameterList.ParameterValueStruct = []

        param = models.ParameterValueStruct()
        param.Name = 'InternetGatewayDevice.ManagementServer.PeriodicInformEnable'
        param.Value = models.anySimpleType(type='xsd:boolean')
        param.Value.Data = '1'
        request.ParameterList.ParameterValueStruct.append(param)

        param = models.ParameterValueStruct()
        param.Name = 'InternetGatewayDevice.ManagementServer.ConnectionRequestUsername'
        param.Value = models.anySimpleType(type='xsd:string')
        param.Value.Data = '00147F-SpeedTouch780-CP0611JTLNW'
        request.ParameterList.ParameterValueStruct.append(param)

        param = models.ParameterValueStruct()
        param.Name = 'InternetGatewayDevice.ManagementServer.ConnectionRequestPassword'
        param.Value = models.anySimpleType(type='xsd:string')
        param.Value.Data = '98ff55fb377bf724c625f60dec448646'
        request.ParameterList.ParameterValueStruct.append(param)

        param = models.ParameterValueStruct()
        param.Name = 'InternetGatewayDevice.ManagementServer.PeriodicInformInterval'
        param.Value = models.anySimpleType(type='xsd:unsignedInt')
        param.Value.Data = '60'
        request.ParameterList.ParameterValueStruct.append(param)

        request.ParameterKey = 'null'

        def side_effect(*_args, **_kwargs):
            return request

        AutoConfigServer.state_machine().handle_tr069_message = Mock(
            side_effect=side_effect)

        self.p.stop()
        self.p = patch.object(self.acs_state_machine, 'handle_tr069_message',
                              Mock(side_effect=side_effect))
        self.p.start()

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [b'']
        ctx, = server.generate_contexts(ctx)

        server.get_in_object(ctx)
        if ctx.in_error is not None:
            raise ctx.in_error

        server.get_out_object(ctx)
        if ctx.out_error is not None:
            raise ctx.out_error

        server.get_out_string(ctx)

        xml_tree = XmlTree()
        match = xml_tree.xml_compare(
            xml_tree.convert_string_to_tree(b''.join(ctx.out_string)),
            xml_tree.convert_string_to_tree(expected_acs_string))
        self.assertTrue(match)

    def test_parse_fault_response(self):
        """ Tests that a fault response from CPE is correctly parsed. """
        # Example CPE->ACS fault response. Copied from:
        # http://djuro82.blogspot.com/2011/05/tr-069-cpe-provisioning.html
        cpe_string = b'''
        <soapenv:Envelope soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
         <soapenv:Header>
        <cwmp:ID soapenv:mustUnderstand="1">1031422463</cwmp:ID>
         </soapenv:Header>
         <soapenv:Body>
           <soapenv:Fault>
            <faultcode>Client</faultcode>
            <faultstring>CWMP fault</faultstring>
            <detail>
             <cwmp:Fault>
              <FaultCode>9003</FaultCode>
              <FaultString>Invalid arguments</FaultString>
              <SetParameterValuesFault>
               <ParameterName>InternetGatewayDevice.WANDevice.1.WANConnectionDevice.3.WANPPPConnection.1.Password</ParameterName>
               <FaultCode>9003</FaultCode>
               <FaultString>Invalid arguments</FaultString>
              </SetParameterValuesFault>
              <SetParameterValuesFault>
               <ParameterName>InternetGatewayDevice.WANDevice.1.WANConnectionDevice.3.WANPPPConnection.1.Username</ParameterName>
               <FaultCode>9003</FaultCode>
               <FaultString>Invalid arguments</FaultString>
              </SetParameterValuesFault>
             </cwmp:Fault>
            </detail>
           </soapenv:Fault>
         </soapenv:Body>
        </soapenv:Envelope>
        '''
        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)
        server.get_in_object(ctx)
        if ctx.in_error is not None:
            raise ctx.in_error

        # Calls function to receive and process message
        server.get_out_object(ctx)

        output_msg = ctx.out_object[0]
        self.assertEqual(type(output_msg), models.Fault)
        self.assertEqual(output_msg.FaultCode, 9003)
        self.assertEqual(output_msg.FaultString, 'Invalid arguments')
        self.assertEqual(
            output_msg.SetParameterValuesFault[1].ParameterName,
            'InternetGatewayDevice.WANDevice.1.WANConnectionDevice.3.WANPPPConnection.1.Username')
        self.assertEqual(output_msg.SetParameterValuesFault[1].FaultCode, 9003)
        self.assertEqual(
            output_msg.SetParameterValuesFault[1].FaultString,
            'Invalid arguments')

    def test_parse_hex_values(self):
        """
        Test that non-utf-8 hex values can be parsed without error
        """
        # Example TR-069 CPE->ACS RPC call. Copied from:
        # http://djuro82.blogspot.com/2011/05/tr-069-cpe-provisioning.html
        cpe_string = b'''
            <soapenv:Envelope soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
                <soapenv:Header>
                    <cwmp:ID soapenv:mustUnderstand="1">0_THOM_TR69_ID</cwmp:ID>
                </soapenv:Header>
                <soapenv:Body>
                    <cwmp:Inform>
                        <DeviceId>
                            <Manufacturer>THOMSON</Manufacturer>
                            <OUI>00147F</OUI>
                            <ProductClass>SpeedTouch 780</ProductClass>
                            <SerialNumber>CP0611JTLNW</SerialNumber>
                        </DeviceId>
                        <Event soap:arrayType="cwmp:EventStruct[04]">
                            <EventStruct>
                                <EventCode>0 BOOTSTRAP</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>1 BOOT</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>2 PERIODIC</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                            <EventStruct>
                                <EventCode>4 VALUE CHANGE</EventCode>
                                <CommandKey></CommandKey>
                            </EventStruct>
                        </Event>
                        <MaxEnvelopes>2</MaxEnvelopes>
                        <CurrentTime>1970-01-01T00:01:09Z</CurrentTime>
                        <RetryCount>05</RetryCount>
                        <ParameterList soap:arrayType="cwmp:ParameterValueStruct[12]">
                            <ParameterValueStruct>
                                <Name>InternetGatewayDevice.DeviceSummary</Name>
                                <Value xsi:type="xsd:string">
                                    \xff\xff\xff\xff\xff</Value>
                            </ParameterValueStruct>
                        </ParameterList>
                    </cwmp:Inform>
                </soapenv:Body>
            </soapenv:Envelope>
            '''

        server = ServerBase(self.app)

        ctx = MethodContext(server, MethodContext.SERVER)
        ctx.in_string = [cpe_string]
        ctx, = server.generate_contexts(ctx)

        if ctx.in_error is not None:
            print('In error: %s' % ctx.in_error)
        self.assertEqual(ctx.in_error, None)

        server.get_in_object(ctx)


class XmlTree():

    @staticmethod
    def convert_string_to_tree(xmlString):

        return ET.fromstring(xmlString)

    def xml_compare(self, x1, x2, excludes=None):
        """
        Compares two xml etrees
        :param x1: the first tree
        :param x2: the second tree
        :param excludes: list of string of attributes to exclude from comparison
        :return:
            True if both files match
        """
        excludes = [] if excludes is None else excludes

        if x1.tag != x2.tag:
            print('Tags do not match: %s and %s' % (x1.tag, x2.tag))
            return False
        for name, value in x1.attrib.items():
            if name not in excludes:
                if x2.attrib.get(name) != value:
                    print('Attributes do not match: %s=%r, %s=%r'
                          % (name, value, name, x2.attrib.get(name)))
                    return False
        for name in x2.attrib.keys():
            if name not in excludes:
                if name not in x1.attrib:
                    print('x2 has an attribute x1 is missing: %s'
                          % name)
                    return False
        if not self.text_compare(x1.text, x2.text):
            print('text: %r != %r' % (x1.text, x2.text))
            return False
        if not self.text_compare(x1.tail, x2.tail):
            print('tail: %r != %r' % (x1.tail, x2.tail))
            return False
        cl1 = x1.getchildren()
        cl2 = x2.getchildren()
        if len(cl1) != len(cl2):
            print('children length differs, %i != %i'
                  % (len(cl1), len(cl2)))
            return False
        i = 0
        for c1, c2 in zip(cl1, cl2):
            i += 1
            if c1.tag not in excludes:
                if not self.xml_compare(c1, c2, excludes):
                    print('children %i do not match: %s'
                          % (i, c1.tag))
                    return False
        return True

    def text_compare(self, t1, t2):
        """
        Compare two text strings
        :param t1: text one
        :param t2: text two
        :return:
            True if a match
        """
        if not t1 and not t2:
            return True
        if t1 == '*' or t2 == '*':
            return True
        return (t1 or '').strip() == (t2 or '').strip()
