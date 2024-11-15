# Store Management Processor System

A robust service designed to efficiently process thousands of store visit images, ensuring accurate management and streamlined operations for retail businesses. Additionally, I created a frontend interface for this service, which can be accessed at [store-management-processor.fly.dev](https://store-management-processor.fly.dev/)

## Project Overview

Retail Pulse requires a system to handle and process large volumes of images collected from various stores. This application addresses that need by:

1. **Receiving Jobs with Image URLs and Store IDs**:
   - Accepts multiple concurrent jobs, each containing numerous images.
   - Jobs may take from a few minutes to an hour to complete.
  
![Store Management Processor](https://github.com/user-attachments/assets/59d972f3-d30b-4f3a-beef-114acdb5ba32)

2. **Processing Images**:
   - Downloads each image and calculates its perimeter using the formula `2 * (Height + Width)`.
   - Introduces a random sleep time between 0.1 to 0.4 seconds to simulate GPU processing.

3. **Validating Store Information**:
   - Cross-references submitted `store_id` with the provided [Store Master](https://drive.google.com/file/d/1dCdAFEBzN1LVUUKxIZyewOeYx42PtEzb/view?usp=sharing) data, which includes `store_id`, `store_name`, and `area_code`.

4. **Job Tracking**:
   - Provides APIs to submit jobs and retrieve their status and results.

[![Store Analytics Processor](https://github.com/user-attachments/assets/0ef772f6-e37d-4abc-a060-0197509ce5ad)](https://store-management-processor.fly.dev/)

## Features

- **High Throughput Processing**: Efficiently manages multiple jobs with thousands of images.
- **Accurate Store Validation**: Ensures all `store_id`s are validated against the master data.
- **Simulated GPU Processing**: Mimics real-world processing delays to provide realistic performance metrics.
- **Comprehensive Job Tracking**: Allows users to monitor the status and results of their submitted jobs.

## Assumptions

- **Store Master Data**: Assumes the provided CSV file contains accurate and up-to-date store information.
- **Image Processing**: Focuses solely on perimeter calculation without additional image transformations.
- **Simulated Processing Delay**: The random sleep time is used to emulate GPU processing times.

## Deployment

The application is deployed and accessible at: [Store Management Processor System](https://store-management-processor.fly.dev/)

### API Endpoints

1. **Submit Job**
   - **Endpoint**: `POST /api/submit`
   - **Description**: Submits a new job for processing store visit images.
   - **Request Body**:

```json
{
  "count": 2,
  "visits": [
    {
      "store_id": "RP00001",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/2.jpg",
        "https://www.gstatic.com/webp/gallery/3.jpg"
      ],
      "visit_time": "2024-03-15T10:00:00Z"
    },
    {
      "store_id": "RP00002",
      "image_url": [
        "https://www.gstatic.com/webp/gallery/4.jpg"
      ],
      "visit_time": "2024-03-15T11:00:00Z"
    }
  ]
}
```

2. **Check Job Status**
   - **Endpoint**: `GET /api/status?jobid=<job_id>`
   - **Description**: Retrieves the status and results of a submitted job.
   - **Parameters**:
     - `jobid`: The unique identifier of the job.

## Installation and Setup

### Prerequisites

- **Go**: Ensure Go is installed on your system.
- **Docker**: Required for containerized setup.

### Using Docker

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/hitaarthh/store-management-processor
   cd store-management-processor
   ```

2. **Build and Run the Docker Container**:
   ```bash
   docker build -t store-management-processor .
   docker run -p 8081:8081 store-management-processor
   ```

3. **Access the Application**:
   Open `http://localhost:8081` in your browser.

### Without Docker

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/hitaarthh/Store-Management-Processor
   cd store-management-processor
   ```

2. **Run the Application**:
   ```bash
   go run main.go
   ```

3. **Access the Application**:
   Open `http://localhost:8081` in your browser.

## Testing Instructions

- **Submit a Job**: Use tools like Postman or curl to send a POST request to `/api/submit` with the required JSON body.
- **Check Job Status**: Send a GET request to `/api/status?jobid=<job_id>` to retrieve the status of a submitted job.
- **Error Handling**: Test with invalid `store_id`s or malformed JSON to ensure proper error responses.

## Work Environment

- **Operating System**: macOS Monterey 12.0.1
- **Text Editor/IDE**: Visual Studio Code 1.62.3
- **Programming Language**: Go 1.17.3
- **Libraries**: Standard Go libraries (`net/http`, `encoding/json`, `image`, etc.)
- **Tools**: Docker 20.10.8, Fly.io CLI 0.0.250


## Future Improvements

Based on the schema of the master store data and the original requirements, here are some **unique and impactful future improvements** that align with the existing project while creatively leveraging the data and functionality:

#### **1. Geographical Analytics**
   - **Heat Map Visualization**:
     - Create a live heat map that shows store density by area code using the `area_code` column from the store master data.
     - This can be used to identify high-density regions for better resource allocation.

   - **Performance Metrics Comparison**:
     - Compare job completion times, error rates, and visit frequencies for stores within the same area code.
     - Highlight underperforming or high-potential regions.

   - **Route Optimization for Store Visits**:
     - Implement an API that calculates the optimal route for visiting multiple stores in a single trip.
     - This can help field teams save time and fuel by organizing visits more efficiently.


#### **2. Store Performance Dashboard**
   - **Visit Frequency Tracking**:
     - Log visit timestamps and calculate metrics like the average number of visits per store or peak visit times.
     - This would help monitor store engagement trends over time.

   - **Image Processing Metrics**:
     - Track metrics like success/failure rates of image downloads and processing times for each store.
     - Provide insights into which stores are consistently problematic.

   - **Custom Store Grouping**:
     - Allow users to define and group stores based on criteria like store names (e.g., "PAN shops") or area codes.
     - Display aggregated data for these groups to identify patterns or outliers.

#### **3. Smart Store Matching**
   - **Similar Store Recommendations**:
     - Use string similarity algorithms to recommend stores with similar naming patterns.
     - Example: A search for "PAN Corner" could also return "Pan Paradise" and "The Pan Stop."

   - **Automatic Store Categorization**:
     - Categorize stores (e.g., "Hotels," "Pan Shops," "Grocery Stores") using keywords from the `store_name` column.
     - This enables tailored processing rules for specific categories.

   - **Category-Specific Rules**:
     - Implement rules like giving higher priority to certain store types (e.g., hotels) or customizing image validation criteria based on the category.

#### **4. Area-Based Insights**
   - **Trends by Area Code**:
     - Analyze area performance trends like increasing/decreasing visit frequency or average processing times.
     - Use this data to identify regions needing more attention or resources.

   - **Custom Reporting**:
     - Generate detailed reports segmented by area codes, showing metrics like job completion rates, visit counts, and error rates.

   - **Load Balancing for Processing**:
     - Dynamically distribute processing workloads based on store density in each area.
     - Regions with high store counts get assigned more resources to handle the load efficiently.

#### **5. Real-Time Enhancements**
   - **Real-Time Error Reporting**:
     - Enhance the `/api/status` endpoint to provide real-time updates about errors encountered during image processing.
     - Users can take immediate corrective actions instead of waiting for the job to fail.

   - **Live Job Monitoring**:
     - Create a dashboard that displays live progress updates for ongoing jobs, showing the percentage of completed images and any encountered errors.

#### **6. Advanced Analytics and Machine Learning**
   - **Predictive Analytics**:
     - Use historical data to predict job completion times, error probabilities, or store performance trends.
     - Example: Flag stores with historically high error rates for preemptive resolution.

   - **Anomaly Detection**:
     - Implement machine learning models to detect anomalies in visit patterns, such as sudden drops in visit frequency or unusual processing times.

#### **7. Scalability Enhancements**
   - **Distributed Processing**:
     - Migrate image processing to a distributed system to handle higher volumes of concurrent jobs.
     - Use tools like Kubernetes for scaling.

   - **Cloud Integration**:
     - Integrate with cloud services for faster image processing and dynamic resource allocation.

