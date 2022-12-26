#include <Arduino.h>
#include <Arduino_MKRGPS.h>
#include <LoRa.h>
#include <SPI.h>

int counter = 0;
int button = 2;
int buttonState;

// the setup function runs once when you press reset or power the board
void setup()
{
  Serial.begin(9600);
  while (!Serial)
  {
    ; // wait for serial port to connect. Needed for native USB port only
  }
  if (!GPS.begin())
  {
    Serial.println("Failed to initialize GPS!");
    while (1)
      ;
  }
  Serial.println("GPS initialized");
  // initialize digital pin LED_BUILTIN as an output.
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, LOW);

  if (!LoRa.begin(915E6))
  {
    Serial.println("Starting LoRa failed!");
    while (1)
      ;
  }
}

// the loop function runs over and over again forever
void loop()
{
  if (GPS.available())
  {
    digitalWrite(LED_BUILTIN, HIGH); // turn the LED on (HIGH is the voltage level)
    // read GPS values
    float latitude = GPS.latitude();
    float longitude = GPS.longitude();
    float altitude = GPS.altitude();
    float speed = GPS.speed();
    int satellites = GPS.satellites();

    // print GPS values
    Serial.print("Location: ");
    Serial.print(latitude, 7);
    Serial.print(", ");
    Serial.println(longitude, 7);

    Serial.print("Altitude: ");
    Serial.print(altitude);
    Serial.println("m");

    Serial.print("Ground speed: ");
    Serial.print(speed);
    Serial.println(" km/h");

    Serial.print("Number of satellites: ");
    Serial.println(satellites);

    Serial.println();
  }

  buttonState = digitalRead(button);

  if (buttonState == HIGH)
  {
    // send packet
    LoRa.beginPacket();
    LoRa.print("button pressed");
    LoRa.endPacket();
    counter++;
    Serial.print("Sending packet: ");
    Serial.println(counter);
    delay(500);
  }
}
