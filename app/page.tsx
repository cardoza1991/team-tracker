'use client'
import React, { useState, useEffect } from 'react';
import { MapPin, Users, CheckCircle, Search, Plus } from 'lucide-react';
import TeamManagement from '@/components/team-management';
import TeamAssignments from '@/components/team-assignments';
import VisitHistory from '@/components/visit-history';
import { Team, Location, Stats } from '@/types';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const Dashboard = () => {
  const [locations, setLocations] = useState<Location[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(null);
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
  const [notes, setNotes] = useState('');
  const [showCheckInForm, setShowCheckInForm] = useState(false);
  const [stats, setStats] = useState<Stats>({
    totalLocations: 213,
    preachedLocations: 0,
    activeTeams: 0,
    teams: []
  });

  const fetchLocations = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/locations');
      if (response.ok) {
        const data = await response.json();
        setLocations(data);
      }
    } catch (error) {
      console.error('Error fetching locations:', error);
    }
  };

  const fetchStatistics = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/statistics');
      if (response.ok) {
        const data = await response.json();
        setStats(prev => ({
          ...prev,
          totalLocations: data.total_locations,
          preachedLocations: data.preached_locations,
        }));
      }
    } catch (error) {
      console.error('Error fetching statistics:', error);
    }
  };

  useEffect(() => {
    const initializeDashboard = async () => {
      await Promise.all([
        fetchLocations(),
        fetchStatistics()
      ]);
    };

    initializeDashboard();
  }, []);

  const handleCheckIn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedLocation || !selectedTeam) return;

    try {
      const response = await fetch('http://localhost:8080/api/visits', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          location_id: selectedLocation.id,
          team_id: selectedTeam.id,
          notes: notes,
        }),
      });

      if (response.ok) {
        setShowCheckInForm(false);
        setNotes('');
        setSelectedLocation(null);
        setSelectedTeam(null);
        
        await Promise.all([
          fetchStatistics(),
          fetchLocations()
        ]);
      }
    } catch (error) {
      console.error('Error checking in:', error);
    }
  };

  const handleTeamsChange = (teams: Team[]) => {
    setStats(prev => ({
      ...prev,
      activeTeams: teams?.length || 0,
      teams: teams || []
    }));
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Stats Bar */}
      <div className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-3">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card className="bg-blue-50">
              <CardContent className="p-4 flex items-center space-x-4">
                <MapPin className="w-8 h-8 text-blue-600" />
                <div>
                  <p className="text-sm text-gray-600">Locations Preached</p>
                  <p className="text-2xl font-bold text-blue-600">{stats.preachedLocations}/{stats.totalLocations}</p>
                </div>
              </CardContent>
            </Card>
            
            <Card className="bg-green-50">
              <CardContent className="p-4 flex items-center space-x-4">
                <Users className="w-8 h-8 text-green-600" />
                <div>
                  <p className="text-sm text-gray-600">Teams Active</p>
                  <p className="text-2xl font-bold text-green-600">{stats.activeTeams}</p>
                </div>
              </CardContent>
            </Card>

            <Card className="bg-purple-50">
              <CardContent className="p-4 flex items-center space-x-4">
                <CheckCircle className="w-8 h-8 text-purple-600" />
                <div>
                  <p className="text-sm text-gray-600">Progress</p>
                  <p className="text-2xl font-bold text-purple-600">
                    {Math.round((stats.preachedLocations / stats.totalLocations) * 100)}%
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 py-6">
        {/* Search and Action Bar */}
        <div className="mb-6 flex flex-col md:flex-row gap-4">
          <div className="relative flex-grow">
            <Search className="absolute left-3 top-3 h-5 w-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search locations or teams..."
              className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <button
            onClick={() => setShowCheckInForm(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg flex items-center justify-center"
          >
            <Plus className="h-5 w-5 mr-2" />
            Check In Location
          </button>
        </div>

        {/* Main Grid */}
        <div className="grid grid-cols-1 gap-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Teams Panel */}
            <div className="md:col-span-1 space-y-6">
              <TeamManagement onTeamsChange={handleTeamsChange} />
              {selectedTeam && <TeamAssignments teamId={selectedTeam.id} />}
            </div>

            {/* Map Panel */}
            <Card className="md:col-span-2">
              <CardHeader>
                <CardTitle>Hampton Roads Preaching Locations</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-[600px] w-full rounded-lg overflow-hidden">
                  <iframe 
                    src="https://www.google.com/maps/d/embed?mid=1Z-uxGE8zPS6GutCmCBbVtemgTSWGiXg&ehbc=2E312F" 
                    className="w-full h-full border-0"
                    title="Hampton Roads Preaching Locations"
                  />
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Visit History */}
          <div className="col-span-full">
            <VisitHistory />
          </div>
        </div>
      </div>

      {/* Check-in Modal */}
      <Dialog open={showCheckInForm} onOpenChange={setShowCheckInForm}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Check In Location</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleCheckIn} className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label>Select Team</Label>
              <Select
                value={selectedTeam?.id?.toString()}
                onValueChange={(value) => {
                  const team = stats.teams.find(t => t.id === parseInt(value));
                  setSelectedTeam(team);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a team..." />
                </SelectTrigger>
                <SelectContent>
                  {stats.teams.map(team => (
                    <SelectItem key={team.id} value={team.id.toString()}>
                      {team.name} - {team.leader}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label>Select Location</Label>
              <Select
                value={selectedLocation?.id?.toString()}
                onValueChange={(value) => {
                  const location = locations.find(l => l.id === parseInt(value));
                  setSelectedLocation(location);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a location..." />
                </SelectTrigger>
                <SelectContent>
                  {locations.map(location => (
                    <SelectItem key={location.id} value={location.id.toString()}>
                      {location.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <Label>Notes (optional)</Label>
              <Textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Add any notes about the visit..."
                className="resize-none"
              />
            </div>

            <div className="flex justify-end gap-3">
              <Button 
                type="button" 
                variant="outline" 
                onClick={() => setShowCheckInForm(false)}
              >
                Cancel
              </Button>
              <Button 
                type="submit"
                disabled={!selectedLocation || !selectedTeam}
              >
                Check In
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default Dashboard;