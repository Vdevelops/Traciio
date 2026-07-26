Berikut rancangan scope AI yang lebih tepat untuk CRM healthcare/pharma kamu. Fokusnya: AI bukan cuma chatbot, tapi **CRM copilot** yang bisa membaca data, membuat/mengubah data, memberi insight, dan prediksi, namun tetap tunduk ke RBAC/scope user login.

**Prinsip Utama**
AI harus selalu bekerja dengan 4 lapisan:

1. **User intent**
   Apa tujuan user: lihat data, tambah data, update, hapus, insight, prediksi, rekomendasi.

2. **Entity/domain**
   Data apa yang diminta: leads, accounts, contacts, deals, tasks, schedules, visit reports, products, product sales, sales performance, monthly targets, route optimization.

3. **RBAC + scope**
   Data yang boleh dilihat:
   - `admin`: global
   - `sales_manager`: team/bawahan dalam scope brick/team
   - `sales`: own data
   - jika tidak punya permission, AI harus bilang tidak punya akses
   - jika punya akses tapi data kosong, AI harus bilang data kosong, bukan “tidak punya akses”

4. **Action policy**
   - Read/insight/prediction boleh langsung.
   - Create/update single record boleh jika permission ada.
   - Delete/bulk update/high impact wajib confirmation.
   - Semua write action harus audit log.

**Scope Kemampuan AI**
| Scope | Contoh User Command | Yang Harus AI Lakukan |
|---|---|---|
| Read Data | “Tampilkan deal saya bulan ini” | Ambil deals sesuai RBAC/scope |
| Create Data | “Buat task follow up RS Hermina besok jam 10” | Resolve account/contact, create task |
| Update Data | “Pindahkan deal Antibiotik Q3 ke proposal” | Cek permission, validasi deal scope, update stage |
| Delete/Archive | “Hapus task duplikat ini” | Minta konfirmasi + reason |
| Insight | “Kenapa penjualan Juli turun?” | Ambil data sales, deals, visit, produk, jelaskan penyebab |
| Prediction | “Prospect mana yang berpotensi deal?” | Hitung score dari BANT, stage, activity, visit, product interest |
| Recommendation | “Apa prioritas saya hari ini?” | Next best action berdasarkan due task, stale deals, hot leads |
| Forecast | “Prediksi closing bulan ini” | Forecast deal weighted revenue dari probability/stage |
| Monitoring | “Deal mana yang stagnan?” | Cari deal tanpa activity/stage movement > N hari |
| Bulk Action | “Buat follow-up untuk semua lead panas” | Tampilkan preview, minta approval |

**Intent Classifier**
AI sebaiknya mengubah semua pesan user menjadi struktur internal seperti ini:

```json
{
  "intent": "predict",
  "entity": "prospect",
  "scope": "rbac_user_scope",
  "filters": {
    "date_range": "bulan_ini",
    "owner": "saya",
    "stage": null,
    "brick": null
  },
  "action": "rank_deal_potential",
  "output_format": "table_with_insight",
  "confirmation_required": false
}
```

**Daftar Intent yang Perlu Didukung**
| Intent | Kata Pemicu |
|---|---|
| `read` | tampilkan, lihat, berikan data, list, cari |
| `create` | buat, tambahkan, jadwalkan, catat, log |
| `update` | ubah, update, pindahkan, reschedule, tandai |
| `delete` | hapus, delete, arsipkan |
| `insight` | kenapa, analisa, ringkas, jelaskan, performa |
| `predict` | prediksi, potensi, kemungkinan, peluang, forecast |
| `recommend` | saran, rekomendasi, prioritas, next action |
| `compare` | bandingkan, vs, dibanding, growth |
| `monitor` | stagnan, overdue, belum follow up, risiko |
| `bulk_action` | semua, massal, seluruh, buatkan untuk semua |

**Contoh Perintah yang Harus Bisa Dipahami**
Read:
- “Tampilkan lead saya bulan ini.”
- “Berikan data penjualan produk Juli dari yang paling banyak terjual.”
- “List deal bawahan saya yang stage negotiation.”

CRUD:
- “Buat task follow up untuk RS Kariadi besok jam 09:00.”
- “Tambahkan activity call untuk lead PT Sehat Farma.”
- “Update status task follow up menjadi completed.”
- “Pindahkan deal Antibiotik Q3 ke closed won.”

Insight:
- “Kenapa target saya belum tercapai?”
- “Kenapa conversion rate bulan ini turun?”
- “Produk apa yang paling sering saya jual bulan ini?”

Prediction:
- “Prospect mana yang paling berpotensi jadi deal?”
- “Deal mana yang kemungkinan gagal closing?”
- “Prediksi revenue closing bulan ini.”

Recommendation:
- “Apa prioritas aktivitas saya hari ini?”
- “Lead mana yang harus saya follow up dulu?”
- “Account mana yang perlu dikunjungi minggu ini?”

**Konteks Data untuk Prediksi Prospect**
Untuk menilai prospect berpotensi deal, AI perlu mengambil sinyal ini:

| Sinyal | Bobot |
|---|---|
| BANT lengkap dan score tinggi | tinggi |
| Lead status qualified/interested | tinggi |
| Ada product interest | tinggi |
| Ada visit report terbaru | sedang-tinggi |
| Ada activity/follow-up aktif | sedang |
| Deal value besar | sedang |
| Stage sudah proposal/negotiation | tinggi |
| Tidak ada follow-up >14 hari | risiko |
| Task overdue | risiko |
| Account punya histori pembelian | tinggi |

Output prediksi harus seperti:
- Prospect/deal
- Probability/confidence
- Alasan
- Risiko
- Next best action

**Format Jawaban AI**
Untuk semua jawaban data/insight, pakai format konsisten:

1. Jawaban inti
2. Tabel data jika ada
3. Insight
4. Rekomendasi aksi
5. Action suggestion

Contoh:
```md
Dari data yang dapat Anda akses, 3 prospect paling berpotensi deal minggu ini adalah:

| Prospect | Probability | Alasan | Next Action |
|---|---:|---|---|
| RS A | 82% | BANT lengkap, visit terbaru, product interest tinggi | Follow up proposal |
| Klinik B | 74% | Activity aktif, stage negotiation | Jadwalkan call |

Insight:
- Prospect dengan visit terbaru punya peluang lebih tinggi.
- 2 deal berisiko karena belum ada follow-up >14 hari.
                                                                                           
Rekomendasi:
- Prioritaskan RS A hari ini.
- Buat task follow-up untuk Klinik B?
```

**Guardrail Wajib**
- Jangan invent data.
- Kalau data kosong, bilang data kosong sesuai scope.
- Kalau tidak punya permission, bilang tidak punya akses.
- Jangan menyamakan “data kosong” dengan “tidak punya akses”.
- Semua create/update/delete harus lewat backend tool gate.
- Bulk update/delete wajib confirmation.
- Prediction wajib menyebut alasan dan confidence.
- Semua data mengikuti `UserContext.GetScopedUserIDs(resource)`.

**Prioritas Implementasi**
1. Rapikan intent classifier untuk read/create/update/insight/predict.
2. Tambahkan context builder per domain: sales, customer, product, analytics, management.
3. Tambahkan prediction engine sederhana berbasis scoring rule dulu.
4. Tambahkan confirmation gate untuk bulk/high impact action.
5. Tambahkan audit log AI action.
6. Tambahkan tests untuk RBAC: admin global, manager team, sales own.