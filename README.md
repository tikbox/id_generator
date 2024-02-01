# IDGenerator

UniqueIDGenerator is a high-performance and reliable Go library designed to generate unique IDs with exceptional efficiency. It offers several distinct advantages over traditional ID generation approaches.

## Key Features

1. **Uniqueness:** The UniqueIDGenerator generates IDs that are guaranteed to be unique, eliminating the risk of duplicates. Each generated ID has a length of 6 characters, providing a balance between uniqueness and compactness.

2. **Randomness:** Unlike common auto-incrementing ID schemes, UniqueIDGenerator generates IDs in a truly random manner. This randomness ensures that the generated IDs are unpredictable and cannot be easily guessed or exploited.

3. **Concurrency Safety:** The UniqueIDGenerator is designed to handle concurrent usage seamlessly. It employs built-in synchronization mechanisms, such as mutexes, to ensure thread safety, eliminating any concerns related to concurrent access.

4. **Exceptional Performance:** UniqueIDGenerator is highly optimized for performance, enabling rapid ID generation even under heavy workloads. Its efficient implementation minimizes computational overhead and delivers lightning-fast results.

5. **Standalone ID Generation Service:** UniqueIDGenerator can serve as a dedicated ID generation service, allowing other business components or systems to request unique IDs. This decoupled architecture promotes modularity and simplifies integration with various applications.

Integrate UniqueIDGenerator into your Go projects effortlessly and enjoy the benefits of unique, non-sequential, concurrent-safe, and high-performance ID generation. Whether you need IDs for short URLs, database records, or any other use case, UniqueIDGenerator has you covered.

## Design Approach

The design of id_generator revolves around the following key elements:

1. Generating Random IDs: The project generates a specified number of random IDs, such as 1 million, and saves them in a file. These IDs are unique and can be used for various purposes.

2. Mapping IDs to Time: The project maps the first 3600 IDs to the seconds within an hour. By default, each ID is associated with the seconds of the current hour. However, this configuration can be adjusted according to specific needs.

3. ID Retrieval: Users can retrieve an ID by providing a timestamp. The project retrieves the ID corresponding to the given timestamp from the map. Once an ID has been used, it will not be reused. If multiple requests are made with the same timestamp, the first requester will receive an available ID, while subsequent requests will return 0.

4. Periodic Update and ID Synchronization: To ensure a continuous supply of available IDs, the project includes functionality to periodically update the IDs and synchronize them with the file. This process involves generating a new set of IDs for the next time period and updating the file accordingly.

## Usage

Follow the steps below to use id_generator:

1. Initialization: Perform the following steps to initialize the project:
   - Adjust the configuration parameters, such as the number of IDs to generate and the mapping strategy for timestamps, according to your specific requirements.
   - Build and run the project to generate the initial set of random IDs and save them to a file.

2. ID Retrieval: To retrieve an ID, follow these instructions:
   - Get the current timestamp in seconds.
   - Pass the timestamp to the project's `GetID` function.
   - If a non-zero ID is returned, it is available for use. If 0 is returned, no available ID could be found for the given timestamp.

3. Periodic Update and ID Synchronization: To ensure a continuous supply of available IDs, follow these steps:
   - Determine the frequency at which you want to update the IDs (e.g., every hour).
   - Implement a mechanism to trigger the update process at the specified frequency.
   - Within the update process, generate a new set of IDs for the next time period and synchronize them with the file using the `SyncIDsToFile` function.

Please note that the design mentioned above uses an hour and a second as the default configuration, but you can modify it according to your specific needs by referring to the code.

## Contributions and Feedback

Contributions to id_generator are highly appreciated. If you encounter any issues, have feature requests, or would like to contribute to the project, please open an issue on the GitHub repository. We welcome your feedback and suggestions to improve the project further.

## License

id_generator is released under the [MIT License](LICENSE).
