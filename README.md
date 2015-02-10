# heatingeffect
This is a project for gophergala.com event.

Authors:
Marcel Hauf

# chillingeffects
Web API Json client for chillingeffects.org

# harvester
Harvests notices from chillingeffects.org and inserts them into a mongodb.

# discovery
Gets the latest notice ID from chillingeffects and the latest harvested notice ID from the database.
Queues up new harvesting jobs on iron.io for notices not yet in the database.

# webserver
Displays the aggregated data in charts.

# common
Config and common database schemas.
