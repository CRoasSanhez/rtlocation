Application to use realtime updating communication with the client using Pusher and Maps Here

Libraries:
- Pusher
- Maps Here
- JWT

Features:
- Geolocation with mongodb and globalsign/mgo
- JWT integration for api requests
- Pusher to create socket sessions and update View in realtime
- Maps to show vehicles from DB updated in reatl time with Pusher
- Create routes according to the nearest vehicle to the user
- Update vehicles view and map while Updating Vehicle coordinates
- Vehicles History

DB:
- Name: rtlocation in mongodb
- Indexes:
   2d in Vehicles and user collections
