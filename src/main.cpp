#include <Arduino.h>
#include <Wire.h>
#include <LiquidCrystal.h>

// button
const int btm = A4;

const int beep = 9;

// encoder
const int esw = A2, ea = A0, eb = A1;

// lcd
const int rs = 7, en = 6, d4 = 5, d5 = 4, d6 = 3, d7 = 2, backlightLed = 8;

LiquidCrystal lcd(rs, en, d4, d5, d6, d7);

boolean backlight = 0;
int currentLine = 0;
boolean clearOnNext = 0;
boolean clearMsg = 1;
char lastRead;
int num = 0;
// overflow
boolean ovf = 0;
boolean lastEa = 1;
long lastReadEncoderTime = 0;

void led(int isOn) {
    if (isOn) {
        digitalWrite(backlightLed, HIGH);
        backlight = 1;
    } else {
        digitalWrite(backlightLed, LOW);
        backlight = 0;
    }
}

void setup() {
    pinMode(btm, INPUT_PULLUP);
    pinMode(esw, INPUT_PULLUP);
    pinMode(ea, INPUT_PULLUP);
    pinMode(eb, INPUT_PULLUP);
    pinMode(backlightLed, OUTPUT);
    lcd.begin(20, 4);
    led(1);
    // lcd.print(" ");
    // lcd.clear();
    lcd.print(" Sparkle    LCD2USB");
    lcd.setCursor(4, 1);
    lcd.print("cupinkie.com");
    lcd.setCursor(0, 2);
    lcd.print("3D Pinter    LCD2004");
    lcd.setCursor(0, 3);
    lcd.print("Check Serial Port...");
    lcd.setCursor(0, 0);
    Serial.begin(115200);
    lastEa = digitalRead(ea);
}

void loop() {
    if (Serial.available() > 0) {
        // something to read
        lastRead = Serial.read();
        if (clearMsg) {
            lcd.clear();
            clearMsg = 0;
        }
        if (lastRead == '\r') {
            currentLine = 0;
            lcd.setCursor(0, currentLine);
            clearOnNext = 0;
            num = 0;
        } else if (lastRead == '\n') {
            currentLine++;
            if (currentLine > 3) {
                currentLine = 0;
            }
            clearOnNext = 1;
            lcd.setCursor(0, currentLine);
            num = 0;
        } else if (lastRead == 172) {
            led(!backlight);
            Serial.read();
        } else {
            if (clearOnNext) {
                if (ovf) {
                    lcd.clear();
                    ovf = 0;
                } else {
                    lcd.print("                    ");
                    lcd.setCursor(0, currentLine);
                }
                clearOnNext = 0;
            }
            num++;
            if (num > 20) {
                currentLine++;
                if (currentLine > 3) {
                    currentLine = 0;
                }
                lcd.setCursor(0, currentLine);
                num = 1;
                ovf = 1;
            }
            lcd.print(lastRead);
        }
    }
    if (!digitalRead(btm)) {
        Serial.print(1);
        while (!digitalRead(btm)) {
            delay(5);
        }
    }
    if (!digitalRead(esw)) {
        Serial.print(0);
        while (!digitalRead(esw)) {
            delay(5);
        }
    }

    if (millis() > lastReadEncoderTime + 5) {
        // 5ms = 200Hz
        boolean readEa = digitalRead(ea);
        boolean readEb = digitalRead(eb);
        if (!readEa && lastEa) {
            // high 2 low
            if (readEb) {
                Serial.print("+");
            } else {
                Serial.print("-");
            }
        }
        lastEa = readEa;
        lastReadEncoderTime = millis();
    }
}