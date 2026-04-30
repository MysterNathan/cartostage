// // components/eleve/ApplicationModal.tsx
// 'use client'
//
// import { useState, useEffect } from 'react'
// import { createApplication, updateApplication, deleteApplication } from '@/lib/studentApi'
// import type { StudentApplication } from '@/types/eleve'
//
// interface ApplicationModalProps {
//     application: StudentApplication | null
//     isOpen: boolean
//     onClose: () => void
//     onSuccess: (application: StudentApplication) => void
//     onDelete: (id: number) => void
//     isNew: boolean
// }
//
// export default function ApplicationModal({
//                                              application,
//                                              isOpen,
//                                              onClose,
//                                              onSuccess,
//                                              onDelete,
//                                              isNew
//                                          }: ApplicationModalProps) {
//     const [formData, setFormData] = useState({
//         stageId: 0,
//         motivationLetter: '',
//         status: 'pending' as const
//     })
//     const [loading, setLoading] = useState(false)
//     const [error, setError] = useState('')
//     const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
//
//     useEffect(() => {
//         if (application) {
//             setFormData({
//                 stageId: application.stageId,
//                 motivationLetter: application.motivationLetter,
//                 status: application.status
//             })
//         } else {
//             setFormData({
//                 stageId: 0,
//                 motivationLetter: '',
//                 status: 'pending'
//             })
//         }
//         setError('')
//     }, [application])
//
//     const handleSubmit = async (e: React.FormEvent) => {
//         e.preventDefault()
//         setLoading(true)
//         setError('')
//
//         try {
//             let result: StudentApplication
//             if (isNew) {
//                 result = await createApplication(formData)
//             } else {
//                 result = await updateApplication(application!.id, formData)
//             }
//             onSuccess(result)
//         } catch (error) {
//             setError(error instanceof Error ? error.message : 'Erreur lors de la sauvegarde')
//         } finally {
//             setLoading(false)
//         }
//     }
//
//     const handleDelete = async () => {
//         if (!application) return
//
//         setLoading(true)
//         try {
//             await deleteApplication(application.id)
//             onDelete(application.id)
//             onClose()
//         } catch (error) {
//             setError(error instanceof Error ? error.message : 'Erreur lors de la suppression')
//         } finally {
//             setLoading(false)
//             setShowDeleteConfirm(false)
//         }
//     }
//
//     if (!isOpen) return null
//
//     return (
//         <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
//             <div className="bg-white rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
//                 <div className="flex justify-between items-center mb-4">
//                     <h2 className="text-xl font-semibold">
//                         {isNew ? 'Nouvelle candidature' : 'Modifier la candidature'}
//                     </h2>
//                     <button
//                         onClick={onClose}
//                         className="text-gray-400 hover:text-gray-600"
//                     >
//                         ×
//                     </button>
//                 </div>
//
//                 {error && (
//                     <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
//                         {error}
//                     </div>
//                 )}
//
//                 {showDeleteConfirm ? (
//                     <div className="text-center">
//                         <p className="mb-4">Êtes-vous sûr de vouloir supprimer cette candidature ?</p>
//                         <div className="flex gap-2 justify-center">
//                             <button
//                                 onClick={handleDelete}
//                                 disabled={loading}
//                                 className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
//                             >
//                                 {loading ? 'Suppression...' : 'Confirmer'}
//                             </button>
//                             <button
//                                 onClick={() => setShowDeleteConfirm(false)}
//                                 className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
//                             >
//                                 Annuler
//                             </button>
//                         </div>
//                     </div>
//                 ) : (
//                     <form onSubmit={handleSubmit}>
//                         <div className="mb-4">
//                             <label className="block text-sm font-medium text-gray-700 mb-2">
//                                 Lettre de motivation *
//                             </label>
//                             <textarea
//                                 value={formData.motivationLetter}
//                                 onChange={(e) => setFormData({ ...formData, motivationLetter: e.target.value })}
//                                 className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
//                                 rows={6}
//                                 required
//                                 placeholder="Décrivez votre motivation pour ce stage-service..."
//                             />
//                         </div>
//
//                         <div className="mb-6">
//                             <label className="block text-sm font-medium text-gray-700 mb-2">
//                                 Statut
//                             </label>
//                             <select
//                                 value={formData.status}
//                                 onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
//                                 className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
//                                 disabled={!isNew} // Le statut ne peut être modifié que par l'tuteur
//                             >
//                                 <option value="pending">En attente</option>
//                                 <option value="accepted">Acceptée</option>
//                                 <option value="rejected">Refusée</option>
//                             </select>
//                         </div>
//
//                         <div className="flex gap-2">
//                             <button
//                                 type="submit"
//                                 disabled={loading}
//                                 className="flex-1 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
//                             >
//                                 {loading ? 'Sauvegarde...' : isNew ? 'Créer' : 'Modifier'}
//                             </button>
//
//                             {!isNew && (
//                                 <button
//                                     type="button"
//                                     onClick={() => setShowDeleteConfirm(true)}
//                                     className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
//                                 >
//                                     Supprimer
//                                 </button>
//                             )}
//
//                             <button
//                                 type="button"
//                                 onClick={onClose}
//                                 className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
//                             >
//                                 Annuler
//                             </button>
//                         </div>
//                     </form>
//                 )}
//             </div>
//         </div>
//     )
// }
