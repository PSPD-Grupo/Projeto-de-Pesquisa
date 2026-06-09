# Etapa 1: build
FROM maven:3.9-eclipse-temurin-21 AS build
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

# Etapa 2: runtime
FROM eclipse-temurin:21-jre
COPY --from=build /app/target/servidor-b.jar app.jar
ENTRYPOINT ["java", "-jar", "/app.jar"]