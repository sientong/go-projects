import { getNavigate } from "./navigation";

export const fetchWithToken = async (url, options = {}, navigate) => {

    try {

      // Get token from storage (adjust based on your app's storage logic)
      const token = localStorage.getItem("token");
  
      // Add Authorization header
      const headers = {
        ...options.headers,
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      };
  
      // Fetch the data
      const response = await fetch(url, { ...options, headers });
  
      console.log(response);

      if (response.status === 401 || response.status === 403) {
        const navigate = getNavigate();
        navigate('/');
        return null;
      }

      // Check if response is ok (status 200-299)
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
  
      // Return parsed data
      return response;
    } catch (error) {
      console.error("Error in fetchWithToken:", error.message);
      throw error; // Re-throw to handle it in the caller
    }
  };