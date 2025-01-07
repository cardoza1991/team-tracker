'use client'
import React, { useState, useEffect } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

const VisitHistory = () => {
  const [visits, setVisits] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchVisitHistory = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/visits/history');
        if (response.ok) {
          const data = await response.json();
          setVisits(data || []); // Ensure we always set an array
          setError(null);
        } else {
          setError('Failed to fetch visit history');
          setVisits([]); // Set empty array on error
        }
      } catch (error) {
        console.error('Error fetching visit history:', error);
        setError('Failed to load visit history');
        setVisits([]); // Set empty array on error
      } finally {
        setIsLoading(false);
      }
    };

    fetchVisitHistory();
  }, []);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Visit History</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-4">Loading visit history...</div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Visit History</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center text-red-600">{error}</div>
        </CardContent>
      </Card>
    );
  }

  if (!visits || visits.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Visit History</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-4">No visits recorded yet.</div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Visit History</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="py-2 px-4 text-left">Date</th>
                <th className="py-2 px-4 text-left">Team</th>
                <th className="py-2 px-4 text-left">Location</th>
                <th className="py-2 px-4 text-left">Status</th>
                <th className="py-2 px-4 text-left">Notes</th>
              </tr>
            </thead>
            <tbody>
              {visits.map((visit) => (
                <tr key={visit.id} className="border-b">
                  <td className="py-2 px-4">
                    {new Date(visit.visit_date).toLocaleString()}
                  </td>
                  <td className="py-2 px-4">{visit.team_name}</td>
                  <td className="py-2 px-4">{visit.location_name}</td>
                  <td className="py-2 px-4">
                    {visit.is_preached ? (
                      <span className="text-green-600 font-medium">Preached</span>
                    ) : (
                      <span className="text-blue-600 font-medium">Visited</span>
                    )}
                  </td>
                  <td className="py-2 px-4">{visit.notes || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
};

export default VisitHistory;