// File implements the dashcam clip-save command.
//
// Tesla's published car_server.proto deliberately omits the dashcam action
// (VehicleAction field 47 is a gap: 46 = Ping, 48 = AutoSeatClimate), so the
// generated carserver package has no DashcamSaveClipAction type. The action is
// real in vehicle firmware, reached on the infotainment domain — verified
// end-to-end against a 2024 vehicle (Tesla returned {"result": true}). We
// reproduce its exact wire encoding (field 47, empty embedded message) by
// injecting it as a protobuf unknown field, then dispatch it through the same
// signed carserver path every other VehicleAction uses.

package vehicle

import (
	"context"

	carserver "github.com/teslamotors/vehicle-command/pkg/protocol/protobuf/carserver"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// DashcamSaveClip copies the last ~10 minutes of the rolling dashcam buffer
// into the vehicle's protected Saved Clips folder, exactly as a tap on the
// in-car Dashcam icon would. No-op-safe: the vehicle ignores it when dashcam
// is unavailable.
func (v *Vehicle) DashcamSaveClip(ctx context.Context) error {
	va := &carserver.VehicleAction{}
	raw := protowire.AppendTag(nil, 47, protowire.BytesType)
	raw = protowire.AppendBytes(raw, nil) // empty DashcamSaveClipAction{}
	va.ProtoReflect().SetUnknown(protoreflect.RawFields(raw))
	return v.executeCarServerAction(ctx,
		&carserver.Action_VehicleAction{VehicleAction: va})
}
