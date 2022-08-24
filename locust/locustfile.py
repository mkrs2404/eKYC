from locust import HttpUser, between, task


class WebsiteUser(HttpUser):
    wait_time = between(5, 15)

    def on_start(self):
        self.client.post("/api/v1/signup", {
            "email": "sid@one2n.in",
            "name": "sid",
            "plan": "basic"
        })
        
    @task
    def index(self):
        self.client.get("/")

    @task
    def swagger(self):
        self.client.get("/swagger/index.html")

    @task
    def face_match(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/facematch", {
            "image1": "",
            "image2": ""
        })

    @task
    def face_match_sync(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/facematch", {
            "image1": "",
            "image2": ""
        })

    @task
    def get_ocr_data(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/get-score", {
            "job-id": ""
        })

    @task
    def get_score(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/get-score", {
            "job-id": ""
        })

    @task
    def ocr(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/ocr-async", {
            "image": ""
        })

    @task
    def ocr_async(self):
        self.client.headers = {'Authorization': ''}
        self.client.post("/api/v1/ocr-async", {
            "image": ""
        })

