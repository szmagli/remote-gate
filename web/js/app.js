const api = axios.create({
    baseURL: 'http://' + window.location.hostname + ':8080/v1/',
    timeout: 5000,
    headers: {
        'X-Custom-Header': 'foobar'
    }
});


var app = new Vue({
    el: '#app',
    data: {
        login: '',
        password: '',
        stream: '',
        logged: false,
        gate: false,
        timer: 0,
        interval: false,
    },
    mounted: function () {
        this.getStatus()
    },
    methods: {
        authorize: function () {
            console.log(this.login)
            console.log(this.password)
            api.post('login', {
                login: this.login,
                password: this.password
            })
            .then(function (response) {
                console.log(response);
                location.reload();
                this.getStatus()
            })
            .catch(function (error) {
                console.log(error);
                location.reload();
            });
        },
        getStatus: function () {
            const self = this;
            api.post('status', {})
            .then(function (response) {
                self.logged = true;
                self.gate = response.data.data.gate;
                self.timer = response.data.data.default;
                if (response.data.data.running) {
                    self.timer = response.data.data.time;
                    self.tmng()
                }
                self.setCamera(response.data.data.camera);
            })
            .catch(function (error) {
                console.log(error);
            });
        },
        setCamera: function (url) {
            console.log(url)
            this.stream = url;
            console.log(this.$refs)
            var player = document.getElementById('test_video');
            var rtspsource = document.createElement("source");
            rtspsource.src = url;
            console.log(player)
            console.log(this.$refs["video"])
            player.appendChild(rtspsource);
            Streamedian.player('test_video', {
                socket: "wss://streamedian.com/ws/"
            });
            var range = document.getElementById('rate');
            var set_live = document.getElementById('to_end');
            var range_out = document.getElementById('rate_res');
            range.addEventListener('input', function () {
                player.playbackRate = range.value;
                range_out.innerHTML = `x${range.value}`;
            });
            set_live.addEventListener('click', function () {
                range.value = 1.0;
                range_out.innerHTML = `live`;
                player.playbackRate = 1;
                player.currentTime = player.buffered.end(0); //player.seekable.end(player.seekable.length - 1);
            });

        },
        Start: function () {
            const self = this;
            api.post('timing', {
                duration: this.timer
            }).then(function (response) {
                self.gate = true;
                self.tmng();
            })
            .catch(function (error) {
                console.log(error);
            });
        },
        Manual: function () {
            const self = this;
            api.post('manual').then(function (response) {
                self.gate = response.data.data.gate;
            })
            .catch(function (error) {
                console.log(error);
            });
        },
        tmng: function() {
            const self = this;
            if(this.interval) {
                clearInterval(this.interval);
                this.interval = false;
            } else {
                this.interval = setInterval(function () {
                    console.log(self.timer)
                    self.timer--;
                    if (self.timer == 0) {
                        clearInterval(self.interval);
                        self.interval = false;
                        self.gate = false;
                    }
                }, 1000);
            }
        }
    }
})

