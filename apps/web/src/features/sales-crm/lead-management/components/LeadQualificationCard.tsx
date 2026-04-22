'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Switch } from '@/components/ui/switch';
import {
  DollarSign,
  User,
  Package,
  Calendar,
  CheckCircle2,
} from 'lucide-react';
import { useLeadQualification } from '../hooks/useLeadQualification';
import { formatCurrency } from '@/lib/utils';
import type { UpdateLeadQualificationRequest } from '../types/qualification';

interface LeadQualificationCardProps {
  leadId: string;
}

export function LeadQualificationCard({ leadId }: LeadQualificationCardProps) {
  const { qualification, isLoading, updateQualification, isUpdating } = useLeadQualification(leadId);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<UpdateLeadQualificationRequest>({});

  // Initialize formData when entering edit mode
  const handleEditClick = () => {
    if (qualification) {
      setFormData({
        budget_target_amount: qualification.budget_target_amount,
        budget_notes: qualification.budget_notes,
        budget_confirmed: qualification.budget_confirmed,

        authority_target_person: qualification.authority_target_person,
        authority_target_role: qualification.authority_target_role,
        authority_notes: qualification.authority_notes,
        authority_confirmed: qualification.authority_confirmed,

        need_priority_level: qualification.need_priority_level,
        need_notes: qualification.need_notes,
        need_confirmed: qualification.need_confirmed,

        timeline_target_date: qualification.timeline_target_date,
        timeline_flexibility: qualification.timeline_flexibility,
        timeline_notes: qualification.timeline_notes,
        timeline_confirmed: qualification.timeline_confirmed,
      });
    }
    setIsEditing(true);
  };

  if (isLoading) {
    return <QualificationSkeleton />;
  }

  if (!qualification) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          Failed to load qualification data
        </CardContent>
      </Card>
    );
  }

  const handleSave = () => {
    const payload = { ...formData };

    // Format date properly for backend RFC3339
    if (payload.timeline_target_date) {
      try {
        payload.timeline_target_date = new Date(payload.timeline_target_date).toISOString();
      } catch {
        delete payload.timeline_target_date;
      }
    }

    updateQualification(payload);
    setIsEditing(false);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'qualified': return 'bg-green-500 hover:bg-green-600';
      case 'warm': return 'bg-yellow-500 hover:bg-yellow-600';
      case 'cold': return 'bg-blue-500 hover:bg-blue-600';
      case 'unqualified': return 'bg-red-500 hover:bg-red-600';
      default: return 'bg-gray-500 hover:bg-gray-600';
    }
  };

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Qualification Checklist (BANT)</CardTitle>
          <div className="flex items-center gap-2">
            <Badge className={getStatusColor(qualification.qualification_status)}>
              {qualification.qualification_status.toUpperCase()}
            </Badge>
            {!isEditing && (
              <Button variant="outline" size="sm" onClick={handleEditClick}>
                Edit
              </Button>
            )}
          </div>
        </div>
        <div className="flex items-center gap-4 mt-2">
          <div className="flex-1">
            <div className="flex justify-between text-sm mb-1">
              <span>Qualification Score</span>
              <span className="font-medium">{qualification.qualification_score}/100</span>
            </div>
            <Progress value={qualification.qualification_score} className="h-2" />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Accordion type="multiple" defaultValue={['budget', 'authority', 'need', 'timeline']} className="space-y-2">
          {/* Budget Section */}
          <AccordionItem value="budget" className="border rounded-lg px-4">
            <AccordionTrigger className="hover:no-underline py-3">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-full ${qualification.bant_progress?.budget?.completed ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                  {qualification.bant_progress?.budget?.completed ? <CheckCircle2 size={18} /> : <DollarSign size={18} />}
                </div>
                <div className="text-left">
                  <span className="font-medium">Budget</span>
                  <p className="text-sm text-muted-foreground">
                    {qualification.budget_target_amount
                      ? formatCurrency(qualification.budget_target_amount)
                      : 'Not specified'}
                  </p>
                </div>
              </div>
            </AccordionTrigger>
            <AccordionContent className="pb-4">
              {isEditing ? (
                <div className="space-y-3">
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Budget Amount</label>
                    <Input
                      type="number"
                      placeholder="e.g. 50000000"
                      defaultValue={qualification.budget_target_amount ? qualification.budget_target_amount / 100 : ''}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        budget_target_amount: e.target.value ? parseInt(e.target.value) * 100 : undefined
                      }))}
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Notes</label>
                    <Textarea
                      placeholder="Budget notes..."
                      defaultValue={qualification.budget_notes}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        budget_notes: e.target.value
                      }))}
                    />
                  </div>
                  <div className="flex items-center gap-2 pt-2">
                    <Switch
                      checked={formData.budget_confirmed ?? qualification.budget_confirmed}
                      onCheckedChange={(checked) => setFormData(p => ({ ...p, budget_confirmed: checked }))}
                    />
                    <label className="text-sm font-medium">Mark Budget as Confirmed</label>
                  </div>
                </div>
              ) : (
                <div className="space-y-2">
                  {qualification.budget_notes && (
                    <p className="text-sm text-muted-foreground">{qualification.budget_notes}</p>
                  )}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>

          {/* Authority Section */}
          <AccordionItem value="authority" className="border rounded-lg px-4">
            <AccordionTrigger className="hover:no-underline py-3">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-full ${qualification.bant_progress?.authority?.completed ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                  {qualification.bant_progress?.authority?.completed ? <CheckCircle2 size={18} /> : <User size={18} />}
                </div>
                <div className="text-left">
                  <span className="font-medium">Authority</span>
                  <p className="text-sm text-muted-foreground">
                    {qualification.authority_target_person || 'Decision maker not identified'}
                  </p>
                </div>
              </div>
            </AccordionTrigger>
            <AccordionContent className="pb-4">
              {isEditing ? (
                <div className="space-y-3">
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Decision Maker Name</label>
                      <Input
                        placeholder="e.g. John Doe"
                        defaultValue={qualification.authority_target_person}
                        onChange={(e) => setFormData(prev => ({
                          ...prev,
                          authority_target_person: e.target.value
                        }))}
                      />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Role/Position</label>
                      <Input
                        placeholder="e.g. CTO / Director"
                        defaultValue={qualification.authority_target_role}
                        onChange={(e) => setFormData(prev => ({
                          ...prev,
                          authority_target_role: e.target.value
                        }))}
                      />
                    </div>
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Notes</label>
                    <Textarea
                      placeholder="Authority notes..."
                      defaultValue={qualification.authority_notes}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        authority_notes: e.target.value
                      }))}
                    />
                  </div>
                  <div className="flex items-center gap-2 pt-2">
                    <Switch
                      checked={formData.authority_confirmed ?? qualification.authority_confirmed}
                      onCheckedChange={(checked) => setFormData(p => ({ ...p, authority_confirmed: checked }))}
                    />
                    <label className="text-sm font-medium">Mark Authority as Confirmed</label>
                  </div>
                </div>
              ) : (
                <div className="space-y-2">
                  {qualification.authority_target_role && (
                    <p className="text-sm"><span className="font-medium">Role:</span> {qualification.authority_target_role}</p>
                  )}
                  {qualification.authority_notes && (
                    <p className="text-sm text-muted-foreground">{qualification.authority_notes}</p>
                  )}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>

          {/* Need Section */}
          <AccordionItem value="need" className="border rounded-lg px-4">
            <AccordionTrigger className="hover:no-underline py-3">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-full ${qualification.bant_progress?.need?.completed ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                  {qualification.bant_progress?.need?.completed ? <CheckCircle2 size={18} /> : <Package size={18} />}
                </div>
                <div className="text-left">
                  <span className="font-medium">Need</span>
                  <p className="text-sm text-muted-foreground">
                    {qualification.need_target_products && qualification.need_target_products.length > 0
                      ? `${qualification.need_target_products.length} products interested`
                      : 'Products not specified'}
                  </p>
                </div>
              </div>
            </AccordionTrigger>
            <AccordionContent className="pb-4">
              {isEditing ? (
                <div className="space-y-3">
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Priority Level</label>
                    <select
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                      defaultValue={qualification.need_priority_level || ""}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        need_priority_level: e.target.value as any
                      }))}
                    >
                      <option value="" disabled>Select Priority...</option>
                      <option value="low">Low</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                      <option value="critical">Critical</option>
                    </select>
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Notes</label>
                    <Textarea
                      placeholder="Need notes..."
                      defaultValue={qualification.need_notes}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        need_notes: e.target.value
                      }))}
                    />
                  </div>
                  <div className="flex items-center gap-2 pt-2">
                    <Switch
                      checked={formData.need_confirmed ?? qualification.need_confirmed}
                      onCheckedChange={(checked) => setFormData(p => ({ ...p, need_confirmed: checked }))}
                    />
                    <label className="text-sm font-medium">Mark Need as Confirmed</label>
                  </div>
                </div>
              ) : (
                <div className="space-y-2">
                  <p className="text-sm"><span className="font-medium">Priority:</span> <Badge variant="outline" className="capitalize">{qualification.need_priority_level}</Badge></p>
                  {qualification.need_target_products && qualification.need_target_products.length > 0 && (
                    <div className="flex flex-wrap gap-2">
                      {qualification.need_target_products.map((product) => (
                        <Badge key={product.product_id} variant="secondary">
                          {product.product_name}
                        </Badge>
                      ))}
                    </div>
                  )}
                  {qualification.need_notes && (
                    <p className="text-sm text-muted-foreground">{qualification.need_notes}</p>
                  )}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>

          {/* Timeline Section */}
          <AccordionItem value="timeline" className="border rounded-lg px-4">
            <AccordionTrigger className="hover:no-underline py-3">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-full ${qualification.bant_progress?.timeline?.completed ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'}`}>
                  {qualification.bant_progress?.timeline?.completed ? <CheckCircle2 size={18} /> : <Calendar size={18} />}
                </div>
                <div className="text-left">
                  <span className="font-medium">Timeline</span>
                  <p className="text-sm text-muted-foreground">
                    {qualification.timeline_target_date
                      ? new Date(qualification.timeline_target_date).toLocaleDateString()
                      : 'No target date set'}
                  </p>
                </div>
              </div>
            </AccordionTrigger>
            <AccordionContent className="pb-4">
              {isEditing ? (
                <div className="space-y-3">
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Target Date</label>
                      <Input
                        type="date"
                        defaultValue={qualification.timeline_target_date?.split('T')[0]}
                        onChange={(e) => setFormData(prev => ({
                          ...prev,
                          timeline_target_date: e.target.value || undefined
                        }))}
                      />
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-muted-foreground">Flexibility</label>
                      <select
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        defaultValue={qualification.timeline_flexibility || ""}
                        onChange={(e) => setFormData(prev => ({
                          ...prev,
                          timeline_flexibility: e.target.value as any
                        }))}
                      >
                        <option value="" disabled>Select Flexibility...</option>
                        <option value="fixed">Fixed</option>
                        <option value="flexible">Flexible</option>
                        <option value="urgent">Urgent</option>
                      </select>
                    </div>
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-muted-foreground">Notes</label>
                    <Textarea
                      placeholder="Timeline notes..."
                      defaultValue={qualification.timeline_notes}
                      onChange={(e) => setFormData(prev => ({
                        ...prev,
                        timeline_notes: e.target.value
                      }))}
                    />
                  </div>
                  <div className="flex items-center gap-2 pt-2">
                    <Switch
                      checked={formData.timeline_confirmed ?? qualification.timeline_confirmed}
                      onCheckedChange={(checked) => setFormData(p => ({ ...p, timeline_confirmed: checked }))}
                    />
                    <label className="text-sm font-medium">Mark Timeline as Confirmed</label>
                  </div>
                </div>
              ) : (
                <div className="space-y-2">
                  <p className="text-sm">
                    <span className="font-medium">Flexibility:</span> <Badge variant="outline" className="capitalize">{qualification.timeline_flexibility}</Badge>
                  </p>
                  {qualification.timeline_notes && (
                    <p className="text-sm text-muted-foreground">{qualification.timeline_notes}</p>
                  )}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        {isEditing && (
          <div className="flex justify-end gap-2 mt-4 pt-4 border-t">
            <Button variant="outline" onClick={() => setIsEditing(false)}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={isUpdating}>
              {isUpdating ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function QualificationSkeleton() {
  return (
    <Card>
      <CardHeader>
        <div className="h-6 w-48 bg-gray-200 rounded animate-pulse" />
        <div className="h-4 w-full bg-gray-200 rounded animate-pulse mt-2" />
      </CardHeader>
      <CardContent className="space-y-4">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-16 bg-gray-200 rounded animate-pulse" />
        ))}
      </CardContent>
    </Card>
  );
}
