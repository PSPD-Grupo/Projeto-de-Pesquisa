
FROM maven:3.9-eclipse-temurin-21 AS build

WORKDIR /app
COPY pom.xml .

RUN mvn dependency:go-offline -q

COPY src ./src
RUN mvn package -q -DskipTests

FROM eclipse-temurin:21-jre-alpine

WORKDIR /app
COPY --from=build /app/target/servidor-b.jar servidor-b.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "servidor-b.jar"]