'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Stage } from '@/types/stage';
import { addStage, updateStage, deleteStage } from '@/lib/stageApi';
import { authApi } from '@/lib/authApi';
import FiliereSelect from '@/components/admin/FiliereSelect';
import { sanitizeStage } from "@/lib/stageUtils";
import { handleApiError, isAuthError } from "@/lib/errorUtils";

interface StageModalProps {
  stage: Stage | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: (stage: Stage) => void;
  onDelete?: (stageId: number) => void;
  isNew?: boolean;
}

export default function StageModal({
                                     stage,
                                     isOpen,
                                     onClose,
                                     onSuccess,
                                     onDelete,
                                     isNew = false
                                   }: StageModalProps) {
  const router = useRouter();

  // État initial par défaut
  const getInitialStage = (): Stage => ({
    id: Date.now(),
    entreprise: '',
    poste: '',
    adresse: '',
    commune: '',
    sector: '',
    filiere: 'CCST',
    lat: -20.8789, // Coordonnées par défaut pour La Réunion
    lng: 55.4481,
    capacity_total: 1,
    capacity_filled: 0,
    placesDisponibles: 1,
    period: '',
    parcours: 'scolaire',
    familleMetiers: '',
    niveauScolaire: '2de'
  });

  const [formData, setFormData] = useState<Stage>(getInitialStage());
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>('');

  // Initialiser le formulaire
  useEffect(() => {
    if (isOpen) {
      if (stage && !isNew) {
        setFormData(stage);
      } else {
        setFormData(getInitialStage());
      }
      setError('');
    }
  }, [stage, isOpen, isNew]);

  // Mettre à jour les places disponibles automatiquement
  useEffect(() => {
    const disponibles = Math.max(0, formData.capacity_total - formData.capacity_filled);
    if (disponibles !== formData.placesDisponibles) {
      setFormData(prev => ({ ...prev, placesDisponibles: disponibles }));
    }
  }, [formData.capacity_total, formData.capacity_filled]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;

    // Conversion des valeurs numériques
    if (['capacity_total', 'capacity_filled', 'lat', 'lng'].includes(name)) {
      const numValue = parseFloat(value) || 0;
      setFormData(prev => ({ ...prev, [name]: numValue }));
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
  };

  const handleSave = async () => {
    if (!formData.entreprise.trim() || !formData.poste.trim()) {
      setError('L\'enterprise et le poste sont obligatoires');
      return;
    }

    if (formData.capacity_total < 1) {
      setError('La capacité totale doit être d\'au moins 1');
      return;
    }

    if (formData.capacity_filled < 0 || formData.capacity_filled > formData.capacity_total) {
      setError('La capacité remplie doit être entre 0 et la capacité totale');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const stageToSave = sanitizeStage(formData);
      let savedStage: Stage;

      if (isNew) {
        savedStage = await addStage(stageToSave);
      } else {
        savedStage = await updateStage(stageToSave);
      }

      onSuccess(savedStage);
    } catch (error) {
      console.error('Erreur lors de la sauvegarde:', error);

      if (isAuthError(error)) {
        authApi.logout();
        router.push('/login');
        return;
      }

      setError(handleApiError(error));
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!stage?.id || isNew) return;

    const confirmMessage = `Êtes-vous sûr de vouloir supprimer le stage "${stage.poste}" chez ${stage.entreprise} ?`;
    if (!confirm(confirmMessage)) return;

    setIsLoading(true);
    setError('');

    try {
      await deleteStage(stage.id);
      onDelete?.(stage.id);
      onClose();
    } catch (error) {
      console.error('Erreur lors de la suppression:', error);

      if (isAuthError(error)) {
        authApi.logout();
        router.push('/login');
        return;
      }

      setError(handleApiError(error));
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
      <div className="fixed inset-0 bg-none bg-opacity-30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
        <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <h2 className="text-xl font-semibold text-gray-900">
              {isNew ? 'Ajouter un nouveau stage' : 'Modifier le stage'}
            </h2>
            <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 text-2xl font-bold"
                disabled={isLoading}
            >
              ×
            </button>
          </div>

          <div className="p-6 space-y-4">
            {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
                  {error}
                </div>
            )}

            {/* Informations de base */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Entreprise *
                </label>
                <input
                    type="text"
                    name="entreprise"
                    value={formData.entreprise}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                    required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Secteur
                </label>
                <input
                    type="text"
                    name="sector"
                    value={formData.sector}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Poste *
              </label>
              <input
                  type="text"
                  name="poste"
                  value={formData.poste}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                  disabled={isLoading}
                  required
              />
            </div>

            {/* Localisation */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Adresse
                </label>
                <input
                    type="text"
                    name="adresse"
                    value={formData.adresse}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Commune
                </label>
                <input
                    type="text"
                    name="commune"
                    value={formData.commune}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>
            </div>

            {/* Coordonnées GPS */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Latitude
                </label>
                <input
                    type="number"
                    name="lat"
                    value={formData.lat}
                    onChange={handleInputChange}
                    step="0.000001"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Longitude
                </label>
                <input
                    type="number"
                    name="lng"
                    value={formData.lng}
                    onChange={handleInputChange}
                    step="0.000001"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>
            </div>

            {/* Filière sur une ligne complète */}
            <div>
              <FiliereSelect
                  value={formData.filiere}
                  onChange={(value) => setFormData(prev => ({ ...prev, filiere: value }))}
                  disabled={isLoading}
              />
            </div>

            {/* Formation et niveau */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Parcours
                </label>
                <select
                    name="parcours"
                    value={formData.parcours}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                >
                  <option value="scolaire">Scolaire</option>
                  <option value="apprentissage">Apprentissage</option>
                  <option value="mixte">Mixte</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Niveau scolaire
                </label>
                <select
                    name="niveauScolaire"
                    value={formData.niveauScolaire}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                >
                  <option value="2de">2de</option>
                  <option value="1re">1re</option>
                  <option value="Tle">Tle</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Famille de métiers
              </label>
              <input
                  type="text"
                  name="familleMetiers"
                  value={formData.familleMetiers}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                  disabled={isLoading}
                  placeholder="Ex: Électricité, Informatique, BTP..."
              />
            </div>

            {/* Capacité */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Capacité totale
                </label>
                <input
                    type="number"
                    name="capacity_total"
                    value={formData.capacity_total}
                    onChange={handleInputChange}
                    min="1"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Places occupées
                </label>
                <input
                    type="number"
                    name="capacity_filled"
                    value={formData.capacity_filled}
                    onChange={handleInputChange}
                    min="0"
                    max={formData.capacity_total}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                    disabled={isLoading}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Places disponibles
                </label>
                <input
                    type="number"
                    value={formData.placesDisponibles}
                    className="w-full px-3 py-2 border border-gray-200 rounded-md bg-gray-50 text-gray-600"
                    disabled
                    readOnly
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Période
              </label>
              <input
                  type="text"
                  name="period"
                  value={formData.period || ''}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
                  disabled={isLoading}
                  placeholder="Ex: Septembre 2024 - Juin 2025"
              />
            </div>
          </div>

          {/* Actions */}
          <div className="px-6 py-4 border-t border-gray-200 flex justify-between">
            <div>
              {!isNew && onDelete && (
                  <button
                      onClick={handleDelete}
                      disabled={isLoading}
                      className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50"
                  >
                    Supprimer
                  </button>
              )}
            </div>

            <div className="flex space-x-3">
              <button
                  onClick={onClose}
                  disabled={isLoading}
                  className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors disabled:opacity-50"
              >
                Annuler
              </button>
              <button
                  onClick={handleSave}
                  disabled={isLoading}
                  className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 flex items-center"
              >
                {isLoading && (
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                )}
                {isLoading ? 'Sauvegarde...' : 'Sauvegarder'}
              </button>
            </div>
          </div>
        </div>
      </div>
  );
}
