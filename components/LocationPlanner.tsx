// Remove unused imports
import React, { useState, useEffect } from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";

// Rest of the component stays the same

const LocationPlanner = ({ teamId, onPlanComplete }) => {
  const [availableLocations, setAvailableLocations] = useState([]);
  const [selectedLocations, setSelectedLocations] = useState([]);

  useEffect(() => {
    fetchAvailableLocations();
  }, []);

  const fetchAvailableLocations = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/locations/available');
      if (response.ok) {
        const data = await response.json();
        setAvailableLocations(data);
      }
    } catch (error) {
      console.error('Error fetching locations:', error);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://localhost:8080/api/teams/${teamId}/plan`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          location_ids: selectedLocations,
          date: new Date().toISOString().split('T')[0]
        }),
      });

      if (response.ok) {
        onPlanComplete();
      }
    } catch (error) {
      console.error('Error planning visits:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label>Select Locations to Visit</Label>
        <Select
          multiple
          value={selectedLocations}
          onChange={(values) => setSelectedLocations(values)}
        >
          <SelectTrigger>
            <SelectValue placeholder="Choose locations..." />
          </SelectTrigger>
          <SelectContent>
            {availableLocations.map(location => (
              <SelectItem key={location.id} value={location.id.toString()}>
                {location.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <Button type="submit">Plan Visits</Button>
    </form>
  );
};

export default LocationPlanner;