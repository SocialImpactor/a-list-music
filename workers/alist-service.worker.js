const version = "0.0.1";
const cache_version = "music_manager_cache:" + version;

console.log("Music Manager Service Worker Initializing..");

self.addEventListener(
    "install",
    () => {
        console.log("MMServiceWorker installed with cache_version", cache_version);
    });

self.addEventListener("cache", function(evt) {
    console.log("")
});

self.addEventListener("fetch", function(evt){
    console.log(evt);
});
