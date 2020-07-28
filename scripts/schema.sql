CREATE TABLE Measurements(
        Id INT NOT NULL AUTO_INCREMENT,
        Url VARCHAR(255) NOT NULL,
        PRIMARY KEY(Id)
);
CREATE TABLE Probes (
        Id INT NOT NULL AUTO_INCREMENT,
        MeasurementId INT NOT NULL,
        Response VARCHAR(255) NOT NULL,
        Duration DECIMAL NOT NULL,
        CreatedAt DATETIME NOT NULL,
        PRIMARY KEY(Id),
        CONSTRAINT `fk_measurement_id` FOREIGN KEY (MeasurementId) REFERENCES Measurements (Id) ON DELETE CASCADE ON UPDATE RESTRICT
);