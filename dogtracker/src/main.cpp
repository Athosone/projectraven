// libraries
// https://docs.arduino.cc/tutorials/mkr-nb-1500/setting-radio-access
// https://docs.arduino.cc/tutorials/mkr-nb-1500/nb-web-client
// https://wheretheravensleep.com/base64/SFRUUEJJTiBpcyBhd2Vzb21l
#include <MKRNB.h>

const char PINNUMBER[] = "1503";
const char server[] = "example.org";

NBClient client;
NB nbAccess;
GPRS gprs;

void setup()
{
  // initialize serial communications and wait for port to open:
  Serial.begin(9600);
  while (!Serial)
  {
    ; // wait for serial port to connect. Needed for native USB port only
  }

  Serial.println("Starting Arduino web client.");
  // connection state
  boolean connected = false;

  // After starting the modem with NB.begin()
  // attach to the GPRS network with the APN, login and password
  while (!connected)
  {
    if ((nbAccess.begin(PINNUMBER) == NB_READY) &&
        (gprs.attachGPRS() == GPRS_READY))
    {
      connected = true;
    }
    else
    {
      Serial.println("Not connected");
      delay(1000);
    }
  }
}

void loop()
{
  Serial.println("Connected");
  delay(1000);
}