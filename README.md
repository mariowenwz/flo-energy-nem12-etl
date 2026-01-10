# Flo Energy – NEM12 Meter Readings ETL

**Overview**

This project implements a lightweight, production-oriented ETL module that processes NEM12 meter data files and generates normalized meter readings suitable for persistence in a relational database.
The solution is designed to handle very large input files efficiently, model domain concepts explicitly, and demonstrate clear design trade-offs under limited time constraints.

## PRD（Product Requirements）
   
1. Problem to Solve
       
    Flo Energy receives a large volume of electricity meter data in NEM12 file format, which contains hierarchical and time-based consumption records.
   
    The problem is to:
   
    * Parse NEM12 files efficiently (including very large files)
    * Extract relevant meter reading information
    * Transform hierarchical interval data into normalized time-series records
    * Generate insert statements for the meter_readings table
      
    The solution should be as close to production-grade as possible, without including infrastructure or deployment concerns.

2. Benefits

    This solution enables:
   
    * Reliable ingestion of raw meter data into downstream systems
    * A normalized, query-friendly representation of time-series consumption data
    * A clear and extensible foundation for future data processing, validation, or analytics pipelines
    * Safe handling of large files without excessive memory usage

3. Features and Scope
   
   **In Scope**
      
    * Streaming processing of NEM12 files (line-by-line)
    * Parsing of:
        ** 200 records (meter context)
        ** 300 records (interval consumption data)
    * Transformation of interval values into timestamped meter readings
    * Generation of SQL insert statements for the meter_readings table
    * Unit tests for core parsing and transformation logic
   
   **Out of Scope**
   
    * Database connectivity or execution of SQL
    * Infrastructure, deployment, or scheduling
    * Validation against the full NEM12 specification beyond required fields
    * Error recovery or retry mechanisms for malformed files
    * Additional NEM12 record types (e.g. 500 for quality information and 900 for file trailers) are explicitly recognized but do not impact the core transformation pipeline. They are intentionally excluded from processing in this implementation, as they do not affect the generation of normalized meter readings. Future iterations could leverage these records for data quality enrichment and ingestion validation.

4. Input & Output
   
    **Input**

    A NEM12 formatted text file
    Relevant records:
   
    * 200 record: meter context (NMI, interval length)
    * 300 record: interval date and consumption values

    **Output**
   
    SQL INSERT statements for the following schema:
   
        create table meter_readings (
          id uuid default gen_random_uuid() not null,
          nmi varchar(10) not null,
          timestamp timestamp not null,
          consumption numeric not null,
          constraint meter_readings_pk primary key (id),
          constraint meter_readings_unique_consumption unique (nmi, timestamp)
        );
    
    Each interval value in the input produces exactly one output row. There are exactly 48 interval values in the 300 record sample data, which is consistent with the 30 mins interval.

5. Assumptions
   
    * Each 300 record belongs to the most recent preceding 200 record
    * Interval length is consistent within a meter context
    * Consumption values are ordered and map sequentially to interval timestamps
    * Duplicate (nmi, timestamp) conflicts are handled downstream (e.g. via database constraints)
    * Input files are well-formed for the relevant record types

## Technical Design

1. High-Level Architecture
   
    * NEM12 File
           
    * Streaming Reader

    * Record Parser  --> identifies 200 / 300 records

    * Meter Context <-- current 200 record

    * Interval Expander  --> expands intervals into timestamps

    * MeterReading  --> normalized domain model

    * SQL Generator

    This pipeline-oriented design allows the system to process arbitrarily large files while keeping memory usage bounded.

2. Domain Modeling (DDD – Lightweight)

    The solution models the domain using three core concepts:
   
    **MeterContext (Context)**
   
    Represents the contextual metadata derived from a 200 record.
    * nmi
    * intervalLength
    This is a stateful parsing context, not a persisted domain entity.
    
    **MeterReading (Fact)**
   
    Represents a normalized, timestamped meter reading.
    * nmi
    * timestamp
    * consumption
    Each MeterReading directly maps to one row in the meter_readings table.
    
    **Domain Relationship Summary**
   
    The 200 record defines the meter context,
    the 300 record provides raw interval facts as inputDTO,
    and meter readings represent normalized time-series data derived by expanding interval records using the meter context.
    The domain consists of only two stable models: the NMI (meter context) and the meter reading. The 300 record is treated as an input DTO of meter reandings rather than a domain entity.

4. Testing Strategy (TDD)
   
    Test-driven development is applied to core transformation logic:

    **Key Test Categories**
   
    Parser Tests

    * Correct parsing of 200 and 300 records
    * Interval Expansion Tests
    * Correct timestamp generation based on interval length
    * Correct mapping of consumption values to timestamps
   
    End-to-End Tests

    * Given a small NEM12 sample
    * When processed through the pipeline
    * Then expected SQL insert statements are generated
    * The focus is on deterministic, domain-level correctness, rather than infrastructure testing.

    Not Covered Cases

    * SQL unique constraint, leave it to DB check
    * Numerical Precision and rounding
    * All fields of NEM12 file validation

5. Technology Stack & Design Choices
       
    Language: Go
   
    * Strong support for streaming I/O
    * Simple concurrency and memory model
    * Fast compilation and minimal runtime dependencies
    * Well-suited for ETL-style pipelines
   
    Parsing Strategy
   
    * Line-by-line file processing
    * No full-file loading into memory
    * Minimal assumptions about file size
   
    Persistence Strategy
   
    * SQL generation only (no database connection)
    * Keeps scope focused and avoids infrastructure concerns
    * Allows easy integration with different downstream systems

6. Rationale for Design Decisions
   
    * Streaming-first design ensures scalability for large files
    * Explicit domain models improve readability and correctness
    * Separation of parsing, transformation, and output enables easy extension
    * Lightweight DDD and TDD provide structure without over-engineering
    * Clear scope boundaries align with time constraints and assessment goals

## Future Improvements
     
Given more time, the following enhancements could be considered:
   
* Support for additional NEM format like NEM13
* Validation and error handling for malformed input
* Batch SQL generation or direct database integration
* Observability (metrics, structured logging)
* Parallel processing of independent meter sections
* Future enhancements could include leveraging 500 records to capture data quality indicators and validating file completeness using 900 record summaries.

**Final Notes**

This implementation intentionally prioritizes clarity, correctness, and extensibility over completeness.
The goal is to demonstrate sound engineering judgment and a production-oriented mindset within a constrained timeframe.



