# peq.nu-backend
A simple hash URL shotener backend using the pigeonhole principle for pontencially create a vendor-free database cluster based on a project [vuo.be](https://github.com/jeielmosi/vuo.be-backend)

## Create a enviroment
Create a Firebase Project and a Firestore database \
Go to `https://console.firebase.google.com/project/{PROJECT-ID}/settings/serviceaccounts/adminsdk` \
And generate new private key, put in `./env` folder with a name `firebase.json` \
You can update the SERVER_PORT on `.env` file, *but make sure it's the same on* `docker-compose.yml` \
## Start your aplication
run `sudo docker compose up --build` on a terminal
