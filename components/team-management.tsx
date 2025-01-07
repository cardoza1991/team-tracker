'use client'
import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2 } from 'lucide-react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const TeamManagement = ({ onTeamsChange }) => {
  const [teams, setTeams] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTeam, setEditingTeam] = useState(null);
  const [formData, setFormData] = useState({ name: '', leader: '' });

  useEffect(() => {
    fetchTeams();
  }, []);

  const fetchTeams = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/teams');
      if (response.ok) {
        const data = await response.json();
        setTeams(data || []); // Ensure we always set an array
        onTeamsChange?.(data || []); // Pass the data to parent, ensure it's an array
      }
    } catch (error) {
      console.error('Error fetching teams:', error);
      setTeams([]); // Set empty array on error
      onTeamsChange?.([]); // Notify parent of empty array
    }
  };

 // Change these lines in team-management.tsx
const handleSubmit = async (e) => {
  e.preventDefault();
  try {
    const url = editingTeam 
      ? `http://localhost:8080/api/teams/${editingTeam.id}`  // Add base URL
      : 'http://localhost:8080/api/teams';  // Add base URL
    
    const method = editingTeam ? 'PUT' : 'POST';
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    });

    if (response.ok) {
      fetchTeams();
      handleCloseModal();
    } else {
      const errorData = await response.text();
      console.error('Server response:', errorData);  // Added error logging
    }
  } catch (error) {
    console.error('Error saving team:', error);
  }
};

const handleDeleteTeam = async (teamId) => {
  if (!confirm('Are you sure you want to delete this team?')) return;
  
  try {
    const response = await fetch(`http://localhost:8080/api/teams/${teamId}`, {  // Add base URL
      method: 'DELETE',
    });

    if (response.ok) {
      fetchTeams();
    }
  } catch (error) {
    console.error('Error deleting team:', error);
  }
};

  const handleOpenModal = (team = null) => {
    setEditingTeam(team);
    setFormData(team || { name: '', leader: '' });
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingTeam(null);
    setFormData({ name: '', leader: '' });
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <CardTitle>Teams</CardTitle>
        <Button onClick={() => handleOpenModal()} className="flex items-center gap-2">
          <Plus className="h-4 w-4" />
          Add Team
        </Button>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {teams && teams.map(team => (
            <div 
              key={team.id} 
              className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100"
            >
              <div>
                <h3 className="font-medium">{team.name}</h3>
                <p className="text-sm text-gray-600">Leader: {team.leader}</p>
                {team.currentLocation && (
                  <div className="mt-2 text-sm text-blue-600">
                    Currently at: {team.currentLocation}
                  </div>
                )}
              </div>
              <div className="flex gap-2">
                <Button 
                  variant="ghost" 
                  size="sm"
                  onClick={() => handleOpenModal(team)}
                >
                  <Edit2 className="h-4 w-4" />
                </Button>
                <Button 
                  variant="ghost" 
                  size="sm"
                  onClick={() => handleDeleteTeam(team.id)}
                  className="text-red-600 hover:text-red-700"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      </CardContent>

      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingTeam ? 'Edit Team' : 'Add New Team'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Team Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Enter team name"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="leader">Team Leader</Label>
              <Input
                id="leader"
                value={formData.leader}
                onChange={(e) => setFormData({ ...formData, leader: e.target.value })}
                placeholder="Enter team leader name"
                required
              />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleCloseModal}>
                Cancel
              </Button>
              <Button type="submit">
                {editingTeam ? 'Update' : 'Create'} Team
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </Card>
  );
};

export default TeamManagement;