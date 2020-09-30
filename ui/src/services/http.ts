import axios from "axios";

const instance = axios.create({
    baseURL: process.env.REACT_APP_API,
    headers: {
      "content-type": "application/json",
      "authorization": "06d2a493-c056-4250-8e7f-4de357ca8cb9",
    },
    responseType: "json"
});

export default instance