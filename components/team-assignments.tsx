import React, { useState, useEffect } from 'react';
import { Plus, Check } from 'lucide-react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

const TeamAssignments = ({ teamId }) => {
  const [assignments, setAssignments] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showAssignDialog, setShowAssignDialog] = useState(false);
  const [availableLocations, setAvailableLocations] = useState([]);
  const [selectedLocations, setSelectedLocations] = useState([]);

  useEffect(() => {
    if (teamId) {
      fetchAssignments();
      fetchAvailableLocations();
    }
  }, [teamId]);

  const fetchAssignments = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/teams/${teamId}/assignments`);
      if (response.ok) {
        const data = await response.json();
        setAssignments(data);
      }
    } catch (error) {
      console.error('Error fetching assignments:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchAvailableLocations = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/locations/available');
      if (response.ok) {
        const data = await response.json();
        setAvailableLocations(data);
      }
    } catch (error) {
      console.error('Error fetching available locations:', error);
    }
  };

  const handleToggleCompletion = async (assignmentId, isCompleted) => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/teams/${teamId}/assignments/${assignmentId}`,
        {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ is_completed: isCompleted }),
        }
      );

      if (response.ok) {
        fetchAssignments();
      }
    } catch (error) {
      console.error('Error updating assignment:', error);
    }
  };

  const handleAssignLocations = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/teams/${teamId}/assignments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          location_ids: selectedLocations.map(id => parseInt(id)),
        }),
      });

      if (response.ok) {
        setShowAssignDialog(false);
        setSelectedLocations([]);
        fetchAssignments();
      }
    } catch (error) {
      console.error('Error assigning locations:', error);
    }
  };

  if (isLoading) {
    return <div>Loading assignments...</div>;
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Assigned Locations</CardTitle>
        <Button size="sm" onClick={() => setShowAssignDialog(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Assign Locations
        </Button>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {assignments.length === 0 ? (
            <div className="text-center text-gray-500 py-4">
              No locations assigned yet
            </div>
          ) : (
            assignments.map((assignment) => (
              <div
                key={assignment.id}
                className="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50"
              >
                <div className="flex items-center space-x-2">
                  <Checkbox
                    checked={assignment.is_completed}
                    onCheckedChange={(checked) =>
                      handleToggleCompletion(assignment.id, checked)
                    }
                  />
                  <span className={assignment.is_completed ? 'line-through text-gray-500' : ''}>
                    {assignment.location_name}
                  </span>
                </div>
                {assignment.is_completed && (
                  <span className="text-sm text-gray-500">
                    {new Date(assignment.completed_date).toLocaleDateString()}
                  </span>
                )}
              </div>
            ))
          )}
        </div>
      </CardContent>

      <Dialog open={showAssignDialog} onOpenChange={setShowAssignDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Assign Locations</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label>Select Locations</Label>
              <Select
                multiple
                value={selectedLocations}
                onValueChange={setSelectedLocations}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Choose locations..." />
                </SelectTrigger>
                <SelectContent>
                  {availableLocations.map((location) => (
                    <SelectItem key={location.id} value={location.id.toString()}>
                      {location.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setShowAssignDialog(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              onClick={handleAssignLocations}
              disabled={selectedLocations.length === 0}
            >
              Assign
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
};

export default TeamAssignments;