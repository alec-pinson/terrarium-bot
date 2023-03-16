# terrarium-bot-v2
Manage a Terrarium. Temperature, humidity, heating, misting, lighting all fully configurable and customisable via YAML.  
  
I currently run this in a Kubernetes cluster (3 Raspberry Pis using [K3S](https://k3s.io/)). I run it along side my other apps:
- [sonoff-lan-api](https://github.com/alec-pinson/sonoff-lan-api) - for controlling sonoff devices over lan
- [meross-lan-api](https://github.com/alec-pinson/meross-lan-api) - for controlling meross devices over lan
- [gpio-to-api](https://github.com/alec-pinson/gpio-to-api) - for controlling gpio and pulling sensor information

## Features
- Fully customisable
- Turn on/off switches
- Pull data from sensors
- Use the data from sensors to trigger an action e.g. switch off/on 
- Send custom alert notifications e.g. humidity low, temp high etc
- Easy to add custom API endpoints that can be called which then run actions (see example below)
- Time based switches, set what should be on during the day, what should be on during the night
- Sunset/sunrise mode, if you want to turn on some lights a few minutes before others in the morning
- Dry run mode using set `DRY_RUN=true` no switches will be turned on or off
- Debug mode for extra output - `DEBUG=true`

## Configuration

[The configuration](cmd/terrarium-bot-v2/configuration.yaml) should be fully customisable and hopefully easy to understand, everything is trigger/action based e.g.

I have a temperature sensor configured:-
```yaml
sensor:
  - id: temperature
    url: http://terrarium-temp-sensor
    jsonPath: Temperature
```
A heating switch configured:-
```yaml
switch:
  - id: heater
    on: http://meross-lan-api/turnOn/heat
    off: http://meross-lan-api/turnOff/heat
```
This heater is controlled by the below trigger, e.g. heating on when below 28 in the day, on when below 22 at night otherwise the heating is off:-
```yaml
trigger:
  - sensor: temperature
    when:
      day:
        below: 28
      night:
        below: 22
    action:
      - heater.on
    else:
      - heater.off
```

Multiple actions can be given to a trigger, e.g. turning lights off for 5 minutes before misting (I like to give my geckos a warning they're about to get sprayed in the face :slightly_smiling_face:):-
```yaml
  - sensor: humidity
    when:
      day:
        below: 61
    action:
      - small_led.off
      - big_led.off
      - uvb.off
      - sleep.5m
      - mister.on
      - sleep.2s
      - small_led.on
      - big_led.on
      - uvb.on
```

Call an action from a URL endpoint, I use my gpio-to-api app to monitor a gpio button, when pressed it can trigger a call to a URL (I press my button before openning the terrarium doors as humidity drops and triggers misting otherwise):-
```yaml
  - endpoint: /prevent-mist
    action:
      - mister.disable.40m
```

## The End
Hopefully the configuration should be pretty self explanitory, if you get stuck or there are any features you think might be missing then feel free to create an issue :slightly_smiling_face:.

## Known Issues
- My crappy sonoff devices don't support getting on/off status over LAN (might replace them at some point) so rather than checking if something is already on/off, I just set on/off anyway. E.g. A switch might be on and I turn it 'on' again. I have added a quick workaround which you can enable if you want `USE_IN_MEMORY_STATUS=true` which will store the status in memory after it has been changed. However, if you turn on or off a switch outside the application, it will not know you have done so.
