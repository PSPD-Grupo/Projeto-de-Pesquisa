import os


class Settings:
    servidor_a_host: str
    servidor_b_host: str

    def __init__(self) -> None:
        self.servidor_a_host = os.getenv("SERVIDOR_A_HOST", "localhost:50051")
        self.servidor_b_host = os.getenv("SERVIDOR_B_HOST", "localhost:50052")


settings = Settings()
