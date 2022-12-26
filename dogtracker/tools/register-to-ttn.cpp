/*
  First Configuration
  This sketch demonstrates the usage of MKR WAN 1300/1310 LoRa module.
  This example code is in the public domain.
*/

#include <MKRWAN.h>

LoRaModem modem;

// Uncomment if using the Murata chip as a module
// LoRaModem modem(Serial1);

String appEui = "0000000000000000";
String appKey = "1D414BB141513CC4EB263960271BB989";
String devAddr = "A8610A31353F8919";
String nwkSKey;
String appSKey;

void setup()
{
  // put your setup code here, to run once:
  Serial.begin(115200);
  while (!Serial)
    ;
  Serial.println("Welcome to MKR WAN 1300/1310 first configuration sketch");
  Serial.println("Register to your favourite LoRa network and we are ready to go!");
  // change this to your regional band (eg. US915, AS923, ...)
  if (!modem.begin(US915))
  {
    Serial.println("Failed to start module");
    while (1)
    {
    }
  };
  Serial.print("Your module version is: ");
  Serial.println(modem.version());
  if (modem.version() != ARDUINO_FW_VERSION)
  {
    Serial.println("Please make sure that the latest modem firmware is installed.");
    Serial.println("To update the firmware upload the 'MKRWANFWUpdate_standalone.ino' sketch.");
  }
  Serial.print("Your device EUI is: ");
  Serial.println(modem.deviceEUI());

  Serial.println("Connecting to TTN using OTAA");

  int connected;
  connected = modem.joinOTAA(appEui, appKey);

  if (!connected)
  {
    Serial.println("Something went wrong; are you indoor? Move near a window and retry");
    delay(1000);
    return;
  }
  Serial.println("You're connected to TTN");

  delay(5000);

  int err;
  modem.setPort(3);
  modem.beginPacket();
  modem.print("HeLoRA world!");
  err = modem.endPacket(true);
  if (err > 0)
  {
    Serial.println("Message sent correctly!");
  }
  else
  {
    Serial.println("Error sending message :(");
  }
}
int connected;

void loop()
{
  if (!connected)
  {
    Serial.println("Connecting to TTN using OTAA");

    connected = modem.joinOTAA(appEui, appKey);

    if (!connected)
    {
      Serial.println("Something went wrong; are you indoor? Move near a window and retry");
      delay(1000);
      return;
    }
    Serial.println("You're connected to TTN");
    int err;
    modem.setPort(3);
    modem.beginPacket();
    modem.print("HeLoRA world!");
    err = modem.endPacket(true);
    if (err > 0)
    {
      Serial.println("Message sent correctly!");
    }
    else
    {
      Serial.println("Error sending message :(");
    }
  }

  while (modem.available())
  {
    Serial.write(modem.read());
  }
  modem.poll();
}